package service

import (
	"go-fiber-starter/app/module/order/repository"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/app/module/order/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Orders) (orders []*response.Order, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (order *response.Order, err error)
	Store(req request.Order) (err error)
	Update(id uint64, req request.Order) (err error)
	Destroy(id uint64) error
}

func Service(Repo repository.IRepository) IService {
	return &service{
		Repo,
	}
}

type service struct {
	Repo repository.IRepository
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

func (_i *service) Store(req request.Order) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.Order) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
