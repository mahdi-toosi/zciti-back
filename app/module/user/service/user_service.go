package service

import (
	"github.com/bangadam/go-fiber-starter/app/module/user/repository"
	"github.com/bangadam/go-fiber-starter/app/module/user/request"
	"github.com/bangadam/go-fiber-starter/app/module/user/response"
	"github.com/bangadam/go-fiber-starter/utils/paginator"
)

type IService interface {
	All(req request.UsersRequest) (articles []*response.User, paging paginator.Pagination, err error)
	Show(id uint64) (article *response.User, err error)
	Store(req request.UserRequest) (err error)
	Update(id uint64, req request.UserRequest) (err error)
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

func (_i *service) All(req request.UsersRequest) (articles []*response.User, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetUsers(req)
	if err != nil {
		return
	}

	for _, result := range results {
		articles = append(articles, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(id uint64) (article *response.User, err error) {
	result, err := _i.Repo.FindOne(id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req request.UserRequest) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.UserRequest) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
