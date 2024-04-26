package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/product/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.ProductsRequest) (products []*schema.Post, paging paginator.Pagination, err error)
	GetOne(businessID uint64, id uint64) (post *schema.Post, err error)
	Create(product []*schema.Product) (err error)
	Update(products []*schema.Product) error
	Delete(businessID uint64, id uint64) error
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.ProductsRequest) (products []*schema.Post, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Post{}).
		Where("business_id = ?", req.BusinessID).
		Where("type = ?", schema.PostTypeProduct)

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

	err = query.
		Preload("Products").
		Preload("Business").
		Preload("Taxonomies").
		Order("created_at desc").Find(&products).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (post *schema.Post, err error) {
	err = _i.DB.Main.
		Preload("Business").
		Preload("Products").
		Preload("Taxonomies").
		Where("business_id = ?", businessID).
		Where("type = ?", schema.PostTypeProduct).
		First(&post, id).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (_i *repo) Create(product []*schema.Product) (err error) {
	err = _i.DB.Main.Create(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (_i *repo) Update(products []*schema.Product) error {
	if err := _i.DB.Main.Transaction(func(tx *gorm.DB) error {
		for _, product := range products {
			if err := tx.Model(&schema.Product{}).
				Where(&schema.Product{ID: product.ID, BusinessID: product.BusinessID}).
				Updates(product).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (_i *repo) Delete(businessID uint64, id uint64) error {
	return _i.DB.Main.Delete(&schema.Product{}, id).Where("business_id = ?", businessID).Error
}
