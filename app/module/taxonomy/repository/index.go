package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/taxonomy/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Taxonomies) (taxonomies []*schema.Taxonomy, paging paginator.Pagination, err error)
	GetOne(id uint64) (taxonomy *schema.Taxonomy, err error)
	Create(taxonomy *schema.Taxonomy) (err error)
	Update(id uint64, taxonomy *schema.Taxonomy) (err error)
	Delete(id uint64) (err error)
	FindUserByMobile(mobile uint64) (taxonomy *schema.Taxonomy, err error)
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
	query := _i.DB.Main.Model(&schema.Taxonomy{})

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

	err = query.Order("created_at desc").Find(&taxonomies).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (taxonomy *schema.Taxonomy, err error) {
	if err := _i.DB.Main.First(&taxonomy, id).Error; err != nil {
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

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Taxonomy{}, id).Error
}

func (_i *repo) FindUserByMobile(mobile uint64) (taxonomy *schema.Taxonomy, err error) {
	if err := _i.DB.Main.Where("mobile = ?", mobile).First(&taxonomy).Error; err != nil {
		return nil, err
	}

	return taxonomy, nil
}
