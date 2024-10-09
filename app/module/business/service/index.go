package service

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/business/repository"
	"go-fiber-starter/app/module/business/request"
	"go-fiber-starter/app/module/business/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Businesses) (businesses []*response.Business, paging paginator.Pagination, err error)
	Show(id uint64, role schema.UserRole) (business *response.Business, err error)
	Store(req request.Business) (err error)
	Update(id uint64, req request.Business) (err error)
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

func (_i *service) Index(req request.Businesses) (businesses []*response.Business, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		businesses = append(businesses, response.FromDomain(result, schema.URUser))
	}

	return
}

func (_i *service) Show(id uint64, role schema.UserRole) (business *response.Business, err error) {
	result, err := _i.Repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	if role == schema.URBusinessOwner {
		return response.FromDomain(result, schema.URBusinessOwner), nil
	}

	return response.FromDomain(result, schema.URUser), nil

}

func (_i *service) Store(req request.Business) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.Business) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
