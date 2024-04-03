package repository

import (
	"github.com/google/uuid"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/asset/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Assets) (assets []*schema.Asset, assetsSize uint64, paging paginator.Pagination, err error)
	GetOne(id uuid.UUID) (asset *schema.Asset, err error)
	Create(asset *schema.Asset) (err error)
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

func (_i *repo) GetAll(req request.Assets) (assets []*schema.Asset, assetsSize uint64, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Asset{}).Where("business_id", req.BusinessID)

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

	err = query.Preload("User").Preload("Business").Order("created_at desc").Find(&assets).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	err = _i.DB.Main.Raw("SELECT SUM(size) FROM assets where business_id = ?", req.BusinessID).Scan(&assetsSize).Error
	if err != nil {
		return
	}

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

func (_i *repo) Delete(id uuid.UUID) error {
	return _i.DB.Main.Delete(&schema.Asset{}, id).Error
}
