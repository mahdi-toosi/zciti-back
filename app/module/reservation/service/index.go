package service

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	oirequest "go-fiber-starter/app/module/orderItem/request"
	prepository "go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/reservation/repository"
	"go-fiber-starter/app/module/reservation/request"
	"go-fiber-starter/app/module/reservation/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Reservations) (reservations []*response.Reservation, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (reservation *response.Reservation, err error)
	Store(req request.Reservation) (err error)
	Update(id uint64, req request.Reservation) (err error)
	Destroy(id uint64) error
	IsReservable(req oirequest.OrderItem, businessID uint64) (err error)
}

func Service(Repo repository.IRepository, pRepo prepository.IRepository) IService {
	return &service{Repo, pRepo}
}

type service struct {
	Repo        repository.IRepository
	ProductRepo prepository.IRepository
}

func (_i *service) Index(req request.Reservations) (reservations []*response.Reservation, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		reservations = append(reservations, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(businessID uint64, id uint64) (article *response.Reservation, err error) {
	result, err := _i.Repo.GetOne(businessID, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (_i *service) Store(req request.Reservation) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.Reservation) (err error) {
	// TODO : check business id permission
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}

func (_i *service) IsReservable(req oirequest.OrderItem, businessID uint64) (err error) {
	if err = _i.Repo.IsReservable(req, businessID); err != nil {
		return err
	}

	variant, err := _i.ProductRepo.GetOneVariant(businessID, req.ProductID)
	if err != nil {
		return err
	}
	if variant.Meta.UniWashMachineStatus == schema.UniWashMachineStatusOFF {
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "در حال حاضر این دستگاه در دسترس نمی باشد. لطفا دستگاه های دیگر را رزرو کنید.",
		}
	}

	return nil
}
