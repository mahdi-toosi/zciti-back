package service

import (
	"go-fiber-starter/app/module/user/repository"
	"go-fiber-starter/app/module/user/request"
	"go-fiber-starter/app/module/user/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Users) (users []*response.User, paging paginator.Pagination, err error)
	Show(id uint64) (user *response.User, err error)
	Store(req request.User) (err error)
	Update(id uint64, req request.User) (err error)
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

func (_i *service) Index(req request.Users) (users []*response.User, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		users = append(users, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(id uint64) (user *response.User, err error) {
	result, err := _i.Repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req request.User) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.User) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
