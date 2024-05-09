package service

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/order/repository"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/app/module/order/response"
	oirepository "go-fiber-starter/app/module/orderItem/repository"
	prepository "go-fiber-starter/app/module/product/repository"
	urequest "go-fiber-starter/app/module/uniwash/request"
	uniService "go-fiber-starter/app/module/uniwash/service"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Orders) (orders []*response.Order, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (order *response.Order, err error)
	Store(req request.Order) (orderID *uint64, err error)
	StoreUniWash(req urequest.StoreUniWash) (err error)
	Update(id uint64, req request.Order) (err error)
	Destroy(id uint64) error
}

func Service(Repo repository.IRepository, ProductRepo prepository.IRepository, OrderItemRepo oirepository.IRepository, uniService uniService.IService) IService {
	return &service{
		Repo, uniService, ProductRepo, OrderItemRepo,
	}
}

type service struct {
	Repo          repository.IRepository
	UniService    uniService.IService
	ProductRepo   prepository.IRepository
	OrderItemRepo oirepository.IRepository
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

func (_i *service) Store(req request.Order) (orderID *uint64, err error) {
	orderID, err = _i.Repo.Create(req.ToDomain())
	if err != nil {
		return nil, err
	}
	return orderID, nil
}

func (_i *service) StoreUniWash(req urequest.StoreUniWash) (err error) {
	if err = _i.UniService.ValidateReservation(req); err != nil {
		return err
	}

	if err = _i.UniService.IsReservable(req); err != nil {
		return err
	}

	// get product
	post, err := _i.ProductRepo.GetOne(req.BusinessID, req.PostID)
	if err != nil {
		return err
	}

	var product schema.Product
	for _, p := range post.Products {
		if p.ID == req.ProductID {
			product = p
			break
		}
	}

	// TODO add these in transactions
	reservationID, err := _i.UniService.Reserve(req)
	if err != nil {
		return err
	}

	// create order
	orderID, err := _i.Repo.Create(&schema.Order{
		ParentID:      nil,
		UserID:        req.UserID,
		TotalAmt:      product.Price,
		BusinessID:    req.BusinessID,
		Meta:          schema.OrderMeta{},
		Status:        schema.OrderStatusPending,
		PaymentMethod: schema.OrderPaymentMethodOnline,
	})
	if err != nil {
		return err
	}

	// create orderItem
	if err := _i.OrderItemRepo.Create(&schema.OrderItem{
		Quantity:      1,
		PostID:        post.ID,
		OrderID:       *orderID,
		ReservationID: reservationID,
		Price:         product.Price,
		Subtotal:      product.Price,
		Type:          schema.OrderItemTypeReservation,
		Meta: schema.OrderItemMeta{
			ProductID:          product.ID,
			ProductTitle:       post.Title,
			ProductType:        product.Type,
			ProductDetail:      product.Meta.Detail,
			ProductVariantType: *product.VariantType,
			//ProductImage:  post.Image,
		},
	}); err != nil {
		return err
	}

	return
}

func (_i *service) Update(id uint64, req request.Order) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
