package service

import (
	"go-fiber-starter/app/module/notificationTemplate/repository"
	"go-fiber-starter/app/module/notificationTemplate/request"
	"go-fiber-starter/app/module/notificationTemplate/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Index) (notificationTemplates []*response.NotificationTemplate, paging paginator.Pagination, err error)
	Store(req request.NotificationTemplate) (err error)
	Update(id uint64, req request.NotificationTemplate) (err error)
	Destroy(businessID uint64, id uint64) error
}

func Service(Repo repository.IRepository) IService {
	return &service{
		Repo,
	}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Index(req request.Index) (notificationTemplates []*response.NotificationTemplate, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		notificationTemplates = append(notificationTemplates, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(id uint64) (notificationTemplate *response.NotificationTemplate, err error) {
	result, err := _i.Repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req request.NotificationTemplate) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.NotificationTemplate) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(businessID uint64, id uint64) error {
	return _i.Repo.Delete(businessID, id)
}
