package repository

import (
	"go-fiber-starter/app/database/schema"
	ntrequest "go-fiber-starter/app/module/notificationTemplate/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req ntrequest.Index) (notificationTemplates []*schema.NotificationTemplate, paging paginator.Pagination, err error)
	GetOne(id uint64) (NotificationTemplate *schema.NotificationTemplate, err error)
	Create(NotificationTemplate *schema.NotificationTemplate) (err error)
	Update(id uint64, NotificationTemplate *schema.NotificationTemplate) (err error)
	Delete(id uint64) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req ntrequest.Index) (notificationTemplates []*schema.NotificationTemplate, paging paginator.Pagination, err error) {
	var total int64

	query := _i.DB.DB.Model(&schema.NotificationTemplate{})
	query.Count(&total)

	req.Pagination.Total = total

	err = query.Offset(req.Pagination.Offset).Limit(req.Pagination.Limit).Find(&notificationTemplates).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (notificationTemplate *schema.NotificationTemplate, err error) {
	if err := _i.DB.DB.First(&notificationTemplate, id).Error; err != nil {
		return nil, err
	}

	return notificationTemplate, nil
}

func (_i *repo) Create(notificationTemplate *schema.NotificationTemplate) (err error) {
	return _i.DB.DB.Create(notificationTemplate).Error
}

func (_i *repo) Update(id uint64, notificationTemplate *schema.NotificationTemplate) (err error) {
	return _i.DB.DB.Model(&schema.NotificationTemplate{}).
		Where(&schema.NotificationTemplate{ID: id}).
		Updates(notificationTemplate).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.DB.Delete(&schema.NotificationTemplate{}, id).Error
}

func (_i *repo) FindUserByMobile(mobile uint64) (notificationTemplate *schema.NotificationTemplate, err error) {
	if err := _i.DB.DB.Where("mobile = ?", mobile).First(&notificationTemplate).Error; err != nil {
		return nil, err
	}

	return notificationTemplate, nil
}
