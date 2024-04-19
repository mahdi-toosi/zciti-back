package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/taxonomy/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Taxonomies) (taxonomies []*schema.Taxonomy, paging paginator.Pagination, err error)
	GetOne(BusinessID uint64, id uint64) (taxonomy *schema.Taxonomy, err error)
	Create(taxonomy *schema.Taxonomy) (err error)
	Update(id uint64, taxonomy *schema.Taxonomy) (err error)
	Delete(BusinessID uint64, id uint64) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Taxonomies) (taxonomies []*schema.Taxonomy, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Debug().Model(&schema.Taxonomy{}).Where("business_id = ?", req.BusinessID)

	if req.Keyword != "" {
		query.Where("title Like ?", "%"+req.Keyword+"%")
	}

	if req.Type != "" {
		query.Where("type = ?", req.Type)
	}

	if req.Domain != "" {
		query.Where("domain = ?", req.Domain)
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at desc").Find(&taxonomies).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (taxonomy *schema.Taxonomy, err error) {
	if err := _i.DB.Main.First(&taxonomy, id).Where("business_id = ?", businessID).Error; err != nil {
		return nil, err
	}

	return taxonomy, nil
}

func (_i *repo) Create(taxonomy *schema.Taxonomy) (err error) {
	return _i.DB.Main.Create(taxonomy).Error
}

func (_i *repo) Update(id uint64, taxonomy *schema.Taxonomy) (err error) {
	return _i.DB.Main.Model(&schema.Taxonomy{}).
		Where(&schema.Taxonomy{ID: id}).
		Updates(taxonomy).Error
}

func (_i *repo) Delete(businessID uint64, id uint64) error {
	return _i.DB.Main.Delete(&schema.Taxonomy{}, id).Where("business_id = ?", businessID).Error
}
