package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/order/repository"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/app/module/order/response"
	oirepository "go-fiber-starter/app/module/orderItem/repository"
	oirequest "go-fiber-starter/app/module/orderItem/request"
	prepository "go-fiber-starter/app/module/product/repository"
	reserveService "go-fiber-starter/app/module/reservation/service"
	uniService "go-fiber-starter/app/module/uniwash/service"
	"go-fiber-starter/internal"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Orders) (orders []*response.Order, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (order *response.Order, err error)
	Store(req request.Order) (orderID *uint64, paymentURL string, err error)
	Status(userID uint64, orderID uint64, authority string) (status string, err error)
	//StoreUniWash(req urequest.StoreUniWash) (err error)
	Update(id uint64, req request.Order) (err error)
	Destroy(id uint64) error
}

func Service(
	config *config.Config,
	repo repository.IRepository,
	zarinPal *internal.ZarinPal,
	uniService uniService.IService,
	productRepo prepository.IRepository,
	orderItemRepo oirepository.IRepository,
	reserveService reserveService.IService,
) IService {
	return &service{
		repo,
		config,
		zarinPal,
		uniService,
		productRepo,
		orderItemRepo,
		reserveService,
	}
}

type service struct {
	Repo           repository.IRepository
	Config         *config.Config
	ZarinPal       *internal.ZarinPal
	UniService     uniService.IService
	ProductRepo    prepository.IRepository
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

func (_i *service) Show(businessID uint64, id uint64) (article *response.Order, err error) {
	result, err := _i.Repo.GetOne(businessID, id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req request.Order) (orderID *uint64, _paymentURL string, err error) {
	// TODO if one of the items was reservation and that is not valid we should not store any thing
	// TODO we should use db transaction ??

	// TODO if order has coupon item with amount of all product prices , status is completed
	req.Status = schema.OrderStatusPending
	// TODO make it dynamic 👇🏻
	req.PaymentMethod = schema.OrderPaymentMethodOnline

	var items []oirequest.ToDomainParams

	for _, item := range req.OrderItems {

		// get product
		post, err := _i.ProductRepo.GetOne(req.BusinessID, item.PostID)
		if err != nil {
			return nil, "", err
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
				return nil, "", err
			}

			if err = _i.ReserveService.IsReservable(item, req.BusinessID); err != nil {
				return nil, "", err
			}

			reservationID, err = _i.UniService.ReserveReservation(item, req.UserID, req.BusinessID)
			if err != nil {
				return nil, "", err
			}
		}

		items = append(items, oirequest.ToDomainParams{
			Post:          *post,
			Product:       product,
			PostID:        item.PostID,
			Quantity:      item.Quantity,
			ReservationID: reservationID,
		})

	}

	var orderItems []schema.OrderItem
	var totalAmt float64
	for _, item := range items {
		i := oirequest.ToDomain(item)
		totalAmt += i.Subtotal
		orderItems = append(orderItems, *i)
	}

	if int(totalAmt) < 100 && int(totalAmt) != 0 {
		return nil, "", &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "ملبغ فاکتور کمتر از حداقل مشخص شده است",
		}
	}

	if int(totalAmt) == 0 {
		req.Status = schema.OrderStatusCompleted
	}

	orderID, err = _i.Repo.Create(req.ToDomain(&totalAmt, nil))
	if err != nil {
		return nil, "", err
	}

	for _, item := range orderItems {
		if err := _i.OrderItemRepo.Create(&item, *orderID); err != nil {
			return nil, "", err
		}
	}

	if int(totalAmt) != 0 {
		callbackURL := fmt.Sprintf("%s/card/payment/result?OrderID=%d&UserID=%d", _i.Config.App.FrontendDomain, *orderID, req.UserID)
		paymentURL, _, _, err := _i.ZarinPal.NewPaymentRequest(int(totalAmt), callbackURL, "رزرو ماشین لباسشویی", "", "")
		if err != nil {
			return nil, "", err
		}

		//order, err := _i.Repo.GetOne(req.BusinessID, *orderID)
		//if err != nil {
		//	return nil, "", err
		//}
		//order.Meta.PaymentAuthority = authority
		//
		//if err := _i.Repo.Update(*orderID, order); err != nil {
		//	return nil, "", err
		//}

		_paymentURL = paymentURL
	}

	return orderID, _paymentURL, nil
}

func (_i *service) Status(userID uint64, orderID uint64, authority string) (status string, err error) {
	// TODO FIX HARD CODE BUSINESS ID !
	order, err := _i.Repo.GetOne(2, orderID)
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

	for _, item := range order.OrderItems {
		if item.Meta.ProductVariantType != schema.ProductVariantTypeWashingMachine {
			continue
		}

		err := _i.UniService.Reserve(*item.ReservationID)
		if err != nil {
			return "OK", err
		}
	}

	return "OK", nil
}

func (_i *service) Update(id uint64, req request.Order) (err error) {
	return _i.Repo.Update(id, req.ToDomain(nil, nil))
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
