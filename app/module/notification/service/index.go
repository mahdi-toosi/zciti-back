package service

import (
	"go-fiber-starter/app/module/notification/repository"
	"go-fiber-starter/app/module/notification/request"
	"go-fiber-starter/app/module/notification/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Notifications) (notifications []*response.Notification, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (notification *response.Notification, err error)
	Store(req request.Notification) (err error)
	Update(id uint64, req request.Notification) (err error)
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

func (_i *service) Index(req request.Notifications) (notifications []*response.Notification, paging paginator.Pagination, err error) {
	notifications, paging, err = _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	return
}

func (_i *service) Show(businessID uint64, id uint64) (article *response.Notification, err error) {
	result, err := _i.Repo.GetOne(businessID, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (_i *service) Store(req request.Notification) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.Notification) (err error) {
	// TODO : check business id permission
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
