package repository

import (
	"github.com/google/uuid"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/asset/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Assets) (assets []*schema.Asset, paging paginator.Pagination, err error)
	GetOne(id uuid.UUID) (asset *schema.Asset, err error)
	Create(asset *schema.Asset) (err error)
	Update(id uuid.UUID, asset *schema.Asset) (err error)
	Delete(id uuid.UUID) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Assets) (assets []*schema.Asset, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Asset{})

	if req.Keyword != "" {
		query.Where("first_name Like ?", "%"+req.Keyword+"%")
		query.Or("last_name Like ?", "%"+req.Keyword+"%")
		query.Or("mobile", req.Keyword)
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Preload("User").Preload("Business").Order("created_at desc").Find(&assets).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uuid.UUID) (asset *schema.Asset, err error) {
	if err := _i.DB.Main.First(&asset, id).Error; err != nil {
		return nil, err
	}

	return asset, nil
}

func (_i *repo) Create(asset *schema.Asset) (err error) {
	return _i.DB.Main.Create(asset).Error
}

func (_i *repo) Update(id uuid.UUID, asset *schema.Asset) (err error) {
	return _i.DB.Main.Model(&schema.Asset{}).
		Where(&schema.Asset{ID: id}).
		Updates(asset).Error
}

func (_i *repo) Delete(id uuid.UUID) error {
	return _i.DB.Main.Delete(&schema.Asset{}, id).Error
}
