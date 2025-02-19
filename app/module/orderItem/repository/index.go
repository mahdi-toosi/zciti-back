package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/orderItem/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.OrderItems) (orderItems []*schema.OrderItem, paging paginator.Pagination, err error)
	GetOne(id uint64) (orderItem *schema.OrderItem, err error)
	Create(orderItem *schema.OrderItem, orderID uint64, tx *gorm.DB) (err error)
	Update(id uint64, orderItem *schema.OrderItem) (err error)
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

func (_i *repo) GetAll(req request.OrderItems) (orderItems []*schema.OrderItem, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.OrderItem{})

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at desc").Find(&orderItems).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (orderItem *schema.OrderItem, err error) {
	if err := _i.DB.Main.
		First(&orderItem, id).
		Error; err != nil {
		return nil, err
	}

	return orderItem, nil
}

func (_i *repo) Create(orderItem *schema.OrderItem, orderID uint64, tx *gorm.DB) (err error) {
	orderItem.OrderID = orderID

	if tx != nil {
		return tx.Create(&orderItem).Error
	}
	return _i.DB.Main.Create(&orderItem).Error
}

func (_i *repo) Update(id uint64, orderItem *schema.OrderItem) (err error) {
	return _i.DB.Main.Model(&schema.OrderItem{}).
		Where(&schema.OrderItem{ID: id}).
		Updates(orderItem).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.OrderItem{}, id).Error
}
