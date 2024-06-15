package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Orders) (orders []*schema.Order, paging paginator.Pagination, err error)
	GetOne(userID uint64, id uint64) (order *schema.Order, err error)
	Create(order *schema.Order) (orderID uint64, err error)
	Update(id uint64, order *schema.Order) (err error)
	Delete(id uint64) (err error)
}

func Repository(db *database.Database) IRepository {
	return &repo{db}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Orders) (orders []*schema.Order, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Order{})

	if req.BusinessID > 0 {
		query.Where(&schema.Order{BusinessID: req.BusinessID})
	}

	if req.UserID > 0 {
		query.Where(&schema.Order{UserID: req.UserID})
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	if req.BusinessID > 0 {
		query.Preload("User")
	}

	err = query.Preload("OrderItems.Reservation").Order("created_at desc").Find(&orders).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(userID uint64, id uint64) (order *schema.Order, err error) {
	if err := _i.DB.Main.
		Where(&schema.Order{UserID: userID}).
		Preload("OrderItems").
		First(&order, id).
		Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (_i *repo) Create(order *schema.Order) (orderID uint64, err error) {
	if err = _i.DB.Main.Create(&order).Error; err != nil {
		return 0, err
	}
	return order.ID, nil
}

func (_i *repo) Update(id uint64, order *schema.Order) (err error) {
	return _i.DB.Main.Model(&schema.Order{}).
		Where(&schema.Order{ID: id, BusinessID: order.BusinessID}).
		Updates(order).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Order{}, id).Error
}
