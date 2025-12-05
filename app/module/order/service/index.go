package service

import (
	"errors"
	"fmt"
	"go-fiber-starter/app/database/schema"
	couponRequst "go-fiber-starter/app/module/coupon/request"
	couponService "go-fiber-starter/app/module/coupon/service"
	"go-fiber-starter/app/module/order/repository"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/app/module/order/response"
	oirepository "go-fiber-starter/app/module/orderItem/repository"
	oirequest "go-fiber-starter/app/module/orderItem/request"
	prepository "go-fiber-starter/app/module/product/repository"
	reserveService "go-fiber-starter/app/module/reservation/service"
	transactionRepo "go-fiber-starter/app/module/transaction/repository"
	uniService "go-fiber-starter/app/module/uniwash/service"
	userService "go-fiber-starter/app/module/user/service"
	walletService "go-fiber-starter/app/module/wallet/service"
	"go-fiber-starter/internal"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/paginator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type IService interface {
	Index(req request.Orders) (orders []*response.Order, paging paginator.Pagination, err error)
	Show(userID uint64, id uint64) (order *response.Order, err error)
	Store(req request.Order) (orderID uint64, paymentURL string, err error)
	Status(userID uint64, orderID uint64, refNum string) (status string, err error)
	Update(id uint64, req request.Order) (err error)
	Destroy(id uint64) error
}

func Service(
	config *config.Config,
	repo repository.IRepository,
	zarinPal *internal.ZarinPal,
	sepGateway *internal.SepGateway,
	uniService uniService.IService,
	userService userService.IService,
	productRepo prepository.IRepository,
	couponService couponService.IService,
	walletService walletService.IService,
	orderItemRepo oirepository.IRepository,
	reserveService reserveService.IService,
	transactionRepo transactionRepo.IRepository,
) IService {
	return &service{
		repo,
		config,
		zarinPal,
		sepGateway,
		uniService,
		userService,
		walletService,
		couponService,
		productRepo,
		reserveService,
		orderItemRepo,
		transactionRepo,
	}
}

type service struct {
	Repo            repository.IRepository
	Config          *config.Config
	ZarinPal        *internal.ZarinPal
	SepGateway      *internal.SepGateway
	UniService      uniService.IService
	UserService     userService.IService
	WalletService   walletService.IService
	CouponService   couponService.IService
	ProductRepo     prepository.IRepository
	ReserveService  reserveService.IService
	OrderItemRepo   oirepository.IRepository
	TransactionRepo transactionRepo.IRepository
}

func (_i *service) Index(req request.Orders) (orders []*response.Order, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		orders = append(orders, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(userID uint64, id uint64) (article *response.Order, err error) {
	result, err := _i.Repo.GetOne(userID, id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req request.Order) (orderID uint64, paymentURL string, err error) {
	// استفاده از تراکنش دیتابیس برای جلوگیری از ثبت ناقص داده‌ها
	tx, err := _i.Repo.BeginTransaction()
	if err != nil {
		return 0, "", err
	}
	defer tx.Rollback() // در صورت بروز خطا، تراکنش لغو شود

	req.Status = schema.OrderStatusPending
	var (
		OrderReservationRanges = make([][]string, 0)
		totalAmt               float64
		orderItems             = make([]schema.OrderItem, 0, len(req.OrderItems))
	)

	for _, item := range req.OrderItems {
		post, err := _i.ProductRepo.GetOne(req.BusinessID, item.PostID)
		if err != nil {
			return 0, "", err
		}

		var product *schema.Product
		for i := range post.Products {
			if post.Products[i].ID == item.ProductID {
				product = &post.Products[i]
				break
			}
		}

		if product == nil {
			return 0, "", errors.New("product not found")
		}

		var reservationID *uint64
		if product.VariantType != nil && *product.VariantType == schema.ProductVariantTypeWashingMachine {
			if err := _i.UniService.ValidateReservation(item); err != nil {
				return 0, "", err
			}
			if err := _i.ReserveService.IsReservable(item, req.BusinessID); err != nil {
				return 0, "", err
			}

			reservationID, err = _i.UniService.ReserveReservation(item, req.User.ID, req.BusinessID)
			if err != nil {
				return 0, "", err
			}

			OrderReservationRanges = append(OrderReservationRanges, []string{item.Date + " " + item.StartTime, item.Date + " " + item.EndTime})
		}

		domainItem := oirequest.ToDomainParams{
			Post:          *post,
			Product:       *product,
			PostID:        item.PostID,
			Quantity:      item.Quantity,
			ReservationID: reservationID,
		}

		i := oirequest.ToDomain(domainItem)
		totalAmt += i.Subtotal
		orderItems = append(orderItems, *i)
	}

	// اعمال کوپن تخفیف در صورت وجود
	if req.CouponCode != "" {
		p := couponRequst.ValidateCoupon{
			OrderTotalAmt:          totalAmt,
			UserID:                 req.User.ID,
			BusinessID:             req.BusinessID,
			Code:                   req.CouponCode,
			OrderReservationRanges: OrderReservationRanges,
		}

		coupon, err := _i.CouponService.ValidateCoupon(p)
		if err != nil {
			return 0, "", err
		}

		if err = _i.CouponService.ApplyCoupon(coupon, req.User.ID, &totalAmt); err != nil {
			return 0, "", err
		}

		req.CouponID = &coupon.ID
	}

	// بررسی مقدار حداقل سفارش
	if totalAmt == 0 {
		req.Status = schema.OrderStatusCompleted
	} else if totalAmt < 100 {
		return 0, "", &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "مبلغ فاکتور کمتر از حداقل مشخص شده است",
		}
	}

	// ایجاد سفارش در دیتابیس
	orderID, err = _i.Repo.Create(req.ToDomain(&totalAmt, nil), tx)
	if err != nil {
		return 0, "", err
	}

	// ایجاد آیتم‌های سفارش
	for _, item := range orderItems {
		if err := _i.OrderItemRepo.Create(&item, orderID, tx); err != nil {
			return 0, "", err
		}
	}

	// در صورت تکمیل شدن سفارش، عملیات پس از ثبت انجام شود
	if req.Status == schema.OrderStatusCompleted {
		if err = _i.UpdateOrderItemsAfterOrderComplete(orderItems); err != nil {
			return orderID, "", nil
		}
	} else {
		redirectURL := fmt.Sprintf(
			"%s/v1/user/orders/status?OrderID=%d&UserID=%d",
			_i.Config.App.BackendDomain,
			orderID,
			req.User.ID,
		)

		var authority string
		if _i.Config.App.Production {
			// paymentURL, authority, _, err = _i.ZarinPal.NewPaymentRequest(
			// 	int(totalAmt), redirectURL,
			// 	"رزرو ماشین لباسشویی",
			// 	"",
			// 	fmt.Sprintf("0%d", req.User.Mobile),
			// )
			paymentURL, err = _i.SepGateway.PaymentService.SendRequest(
				int(totalAmt*1.1)*10,
				strconv.FormatUint(orderID, 10),
				fmt.Sprintf("0%d", req.User.Mobile),
				redirectURL,
			)

			if err != nil {
				return 0, "", err
			}
		} else {
			redirectURL := fmt.Sprintf("%s&Status=%s", redirectURL, "OK")
			paymentURL = redirectURL
			authority = "development authority"
		}

		businessWallet, err := _i.WalletService.GetOrCreateWallet(nil, &req.BusinessID, tx)
		if err != nil {
			return 0, "", err
		}

		err = _i.TransactionRepo.Create(&schema.Transaction{
			Amount:               totalAmt,
			OrderID:              &orderID,
			GatewayTransactionID: &authority,
			UserID:               req.User.ID,
			WalletID:             businessWallet.ID,
			Description:          "رزرو ماشین لباسشویی",
			OrderPaymentMethod:   schema.OrderPaymentMethodOnline,
			Status:               schema.TransactionStatusPending,
		}, tx)

		if err != nil {
			return 0, "", err
		}
	}

	// اضافه کردن کاربر در سیستم
	_ = _i.UserService.InsertUser(req.BusinessID, req.User.ID)

	// تأیید تراکنش دیتابیس
	if err := tx.Commit().Error; err != nil {
		return 0, "", err
	}

	return orderID, paymentURL, nil
}

func (_i *service) Status(userID uint64, orderID uint64, refNum string) (state string, err error) {
	transaction, err := _i.TransactionRepo.GetOne(nil, &orderID)
	if err != nil {
		return "FAILED", err
	}

	order, err := _i.Repo.GetOne(userID, orderID)
	if err != nil {

		return "FAILED", err
	}

	if _i.Config.App.Production {
		// verified, _, _, err := _i.ZarinPal.PaymentVerification(int(order.TotalAmt), authority)
		verified, err := _i.SepGateway.PaymentService.Verify(refNum)
		if err != nil {
			transaction.Status = schema.TransactionStatusFailed
			_ = _i.TransactionRepo.Update(transaction.ID, transaction)

			return "FAILED", err
		}

		if !verified.Success {
			transaction.Status = schema.TransactionStatusFailed
			_ = _i.TransactionRepo.Update(transaction.ID, transaction)

			return "FAILED", &fiber.Error{Code: fiber.StatusBadRequest, Message: "پرداخت ناموفق بوده است"}
		}
	}

	transaction.Status = schema.TransactionStatusSuccess
	err = _i.TransactionRepo.Update(transaction.ID, transaction)
	if err != nil {
		return "", err
	}

	wallet, err := _i.WalletService.Show(&transaction.WalletID, nil, nil)
	if err != nil {
		return "", err
	}
	wallet.Amount += transaction.Amount
	err = _i.WalletService.Update(wallet.ID, schema.Wallet{ID: wallet.ID, Amount: wallet.Amount})
	if err != nil {
		return "", err
	}

	order.Status = schema.OrderStatusCompleted
	if err := _i.Repo.Update(orderID, order); err != nil {
		return "OK", err
	}

	if err = _i.UpdateOrderItemsAfterOrderComplete(order.OrderItems); err != nil {
		return "OK", err
	}

	return "OK", nil
}

func (_i *service) UpdateOrderItemsAfterOrderComplete(orderItems []schema.OrderItem) error {
	for _, item := range orderItems {
		if item.Meta.ProductVariantType == schema.ProductVariantTypeWashingMachine {
			if err := _i.UniService.Reserve(*item.ReservationID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (_i *service) Update(id uint64, req request.Order) (err error) {
	return _i.Repo.Update(id, req.ToDomain(nil, nil))
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
