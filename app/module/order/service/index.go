package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
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
	uniService "go-fiber-starter/app/module/uniwash/service"
	userService "go-fiber-starter/app/module/user/service"
	"go-fiber-starter/internal"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Orders) (orders []*response.Order, paging paginator.Pagination, err error)
	Show(userID uint64, id uint64) (order *response.Order, err error)
	Store(req request.Order) (orderID uint64, paymentURL string, err error)
	Status(userID uint64, orderID uint64, authority string) (status string, err error)
	Update(id uint64, req request.Order) (err error)
	Destroy(id uint64) error
}

func Service(
	config *config.Config,
	repo repository.IRepository,
	zarinPal *internal.ZarinPal,
	uniService uniService.IService,
	userService userService.IService,
	productRepo prepository.IRepository,
	couponService couponService.IService,
	orderItemRepo oirepository.IRepository,
	reserveService reserveService.IService,
) IService {
	return &service{
		repo,
		config,
		zarinPal,
		uniService,
		userService,
		productRepo,
		couponService,
		orderItemRepo,
		reserveService,
	}
}

type service struct {
	Repo           repository.IRepository
	Config         *config.Config
	ZarinPal       *internal.ZarinPal
	UniService     uniService.IService
	UserService    userService.IService
	ProductRepo    prepository.IRepository
	CouponService  couponService.IService
	OrderItemRepo  oirepository.IRepository
	ReserveService reserveService.IService
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
	// TODO if one of the items was reservation and that is not valid we should not store any thing
	// TODO we should use db transaction ??
	req.Status = schema.OrderStatusPending

	var items = make([]oirequest.ToDomainParams, 0)
	var OrderReservationRanges = make([][]string, 0)

	for _, item := range req.OrderItems {
		// get product
		post, err := _i.ProductRepo.GetOne(req.BusinessID, item.PostID)
		if err != nil {
			return 0, "", err
		}

		var product schema.Product
		for _, p := range post.Products {
			if p.ID == item.ProductID {
				product = p
				break
			}
		}

		var reservationID *uint64
		if product.VariantType != nil && *product.VariantType == schema.ProductVariantTypeWashingMachine {
			if err = _i.UniService.ValidateReservation(item); err != nil {
				return 0, "", err
			}

			if err = _i.ReserveService.IsReservable(item, req.BusinessID); err != nil {
				return 0, "", err
			}

			reservationID, err = _i.UniService.ReserveReservation(item, req.User.ID, req.BusinessID)
			if err != nil {
				return 0, "", err
			}

			OrderReservationRanges = append(
				OrderReservationRanges,
				[]string{item.Date + " " + item.StartTime, item.Date + " " + item.EndTime},
			)
		}

		items = append(items, oirequest.ToDomainParams{
			Post:          *post,
			Product:       product,
			PostID:        item.PostID,
			Quantity:      item.Quantity,
			ReservationID: reservationID,
		})
	}

	var totalAmt float64
	var orderItems = make([]schema.OrderItem, 0)

	for _, item := range items {
		i := oirequest.ToDomain(item)
		totalAmt += i.Subtotal
		orderItems = append(orderItems, *i)
	}

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

	if int(totalAmt) == 0 {
		req.Status = schema.OrderStatusCompleted
	} else if int(totalAmt) < 100 {
		return 0, "", &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "ملبغ فاکتور کمتر از حداقل مشخص شده است",
		}
	}
	orderID, err = _i.Repo.Create(req.ToDomain(&totalAmt, nil))
	if err != nil {
		return 0, "", err
	}

	for _, item := range orderItems {
		if err := _i.OrderItemRepo.Create(&item, orderID); err != nil {
			return 0, "", err
		}
	}

	if req.Status == schema.OrderStatusCompleted {
		if err = _i.UpdateOrderItemsAfterOrderComplete(orderItems); err != nil {
			return orderID, "", nil
		}
	} else {
		callbackURL := fmt.Sprintf("%s/v1/user/orders/status?OrderID=%d&UserID=%d",
			_i.Config.App.BackendDomain,
			orderID,
			req.User.ID)

		_paymentURL, _, _, err := _i.ZarinPal.NewPaymentRequest(
			int(totalAmt), callbackURL,
			"رزرو ماشین لباسشویی",
			"",
			fmt.Sprintf("0%d", req.User.Mobile),
		)
		if err != nil {
			return 0, "", err
		}

		paymentURL = _paymentURL
	}

	_ = _i.UserService.InsertUser(req.BusinessID, req.User.ID)

	return orderID, paymentURL, nil
}

func (_i *service) Status(userID uint64, orderID uint64, authority string) (status string, err error) {
	order, err := _i.Repo.GetOne(userID, orderID)
	if err != nil {
		return "NOK", err
	}

	verified, _, _, err := _i.ZarinPal.PaymentVerification(int(order.TotalAmt), authority)
	if err != nil {
		return "NOK", err
	}

	if !verified {
		return "NOK", &fiber.Error{Code: fiber.StatusBadRequest, Message: "پرداخت ناموفق بوده است"}
	}

	// if order.Meta.PaymentAuthority == authority {

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
