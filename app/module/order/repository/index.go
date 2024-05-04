package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Orders) (orders []*schema.Order, paging paginator.Pagination, err error)
	GetOne(businessID uint64, id uint64) (order *schema.Order, err error)
	Create(order *schema.Order) (err error)
	Update(id uint64, order *schema.Order) (err error)
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

func (_i *repo) GetAll(req request.Orders) (orders []*schema.Order, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Order{}).
		Where(&schema.Order{BusinessID: req.BusinessID})

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Preload("User").Order("created_at asc").Find(&orders).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (order *schema.Order, err error) {
	if err := _i.DB.Main.
		Where(&schema.Order{BusinessID: businessID}).
		First(&order, id).
		Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (_i *repo) Create(order *schema.Order) (err error) {
	return _i.DB.Main.Create(order).Error
}

func (_i *repo) Update(id uint64, order *schema.Order) (err error) {
	return _i.DB.Main.Model(&schema.Order{}).
		Where(&schema.Order{ID: id, BusinessID: order.BusinessID}).
		Updates(order).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Order{}, id).Error
}
