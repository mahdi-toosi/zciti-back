package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/business/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Businesses) (businesses []*schema.Business, paging paginator.Pagination, err error)
	GetOne(id uint64) (business *schema.Business, err error)
	Create(business *schema.Business) (err error)
	Update(id uint64, business *schema.Business) (err error)
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

func (_i *repo) GetAll(req request.Businesses) (businesses []*schema.Business, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Business{})

	if req.Keyword != "" {
		query.Where("title Like ?", "%"+req.Keyword+"%")
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	if len(req.IDs) == 0 {
		query.Preload("Owner")
	} else {
		query.Where("id in (?)", req.IDs)
	}

	err = query.Order("created_at desc").Find(&businesses).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (business *schema.Business, err error) {
	if err := _i.DB.Main.Preload("Owner").First(&business, id).Error; err != nil {
		return nil, err
	}

	return business, nil
}

func (_i *repo) Create(business *schema.Business) (err error) {
	return _i.DB.Main.Create(business).Error
}

func (_i *repo) Update(id uint64, business *schema.Business) (err error) {
	return _i.DB.Main.Model(&schema.Business{}).
		Where(&schema.Business{ID: id}).
		Updates(business).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Business{}, id).Error
}
