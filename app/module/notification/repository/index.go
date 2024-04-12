package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/notification/request"
	"go-fiber-starter/app/module/notification/response"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Notifications) (notifications []*response.Notification, paging paginator.Pagination, err error)
	GetOne(businessID uint64, id uint64) (notification *response.Notification, err error)
	Create(notification *schema.Notification) (err error)
	Update(id uint64, notification *schema.Notification) (err error)
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

func (_i *repo) GetAll(req request.Notifications) (notifications []*response.Notification, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Notification{}).
		Where("business_id = ?", req.BusinessID).
		Select("notifications.*, " +
			"users.first_name || ' ' || users.last_name as receiver_full_name, " +
			"businesses.title as business_title").
		Joins("INNER JOIN users ON users.id = notifications.receiver_id").
		Joins("INNER JOIN businesses ON businesses.id = notifications.business_id")

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at asc").Find(&notifications).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (notification *response.Notification, err error) {
	if err := _i.DB.Main.First(&notification, id).Where("business_id = ?", businessID).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

func (_i *repo) Create(notification *schema.Notification) (err error) {
	return _i.DB.Main.Create(notification).Error
}

func (_i *repo) Update(id uint64, notification *schema.Notification) (err error) {
	return _i.DB.Main.Model(&schema.Notification{}).
		Where(&schema.Notification{ID: id, BusinessID: notification.BusinessID}).
		Updates(notification).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Notification{}, id).Error
}
