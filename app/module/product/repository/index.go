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
	GetOneVariant(businessID uint64, id uint64) (product *schema.Product, err error)
	Create(product *schema.Product) (err error)
	Creates(product []*schema.Product) (err error)
	Update(product *schema.Product) (err error)
	CreateAttribute(productID uint64, attrID uint64) error
	DeleteAttribute(productID uint64, attrID uint64) error
	Updates(products []*schema.Product) error
	Delete(businessID uint64, id uint64) error
}

func Repository(db *database.Database) IRepository {
	return &repo{db}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.ProductsRequest) (products []*schema.Post, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Post{}).
		Where(&schema.Post{BusinessID: req.BusinessID, Type: schema.PostTypeProduct})

	if req.CategoryID != "" {
		query = query.
			Joins("JOIN posts_taxonomies ON posts_taxonomies.post_id = posts.id").
			Where("posts_taxonomies.taxonomy_id = ?", req.CategoryID)
	}

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
		Preload("Taxonomies").
		Preload("Products").
		Order("created_at desc").Find(&products).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (post *schema.Post, err error) {
	if err = _i.DB.Main.
		Preload("Taxonomies").
		Preload("Products").
		Where(&schema.Post{BusinessID: businessID, Type: schema.PostTypeProduct}).
		First(&post, id).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (_i *repo) GetOneVariant(businessID uint64, id uint64) (product *schema.Product, err error) {
	if err = _i.DB.Main.
		Where(&schema.Product{BusinessID: businessID}).
		First(&product, id).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (_i *repo) Creates(product []*schema.Product) (err error) {
	err = _i.DB.Main.Create(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (_i *repo) Create(product *schema.Product) (err error) {
	err = _i.DB.Main.Create(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (_i *repo) Update(product *schema.Product) (err error) {
	if err := _i.DB.Main.Model(&schema.Product{}).
		Where(&schema.Product{ID: product.ID, BusinessID: product.BusinessID}).
		Updates(product).Error; err != nil {
		return err
	}
	return nil
}

func (_i *repo) Updates(products []*schema.Product) error {
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

func (_i *repo) CreateAttribute(productID uint64, attrID uint64) (err error) {
	var product schema.Product
	var taxonomy schema.Taxonomy

	if err = _i.DB.Main.First(&product, productID).Error; err != nil {
		return err
	}
	if err = _i.DB.Main.First(&taxonomy, attrID).Error; err != nil {
		return err
	}

	// Remove taxonomy from the product
	if err = _i.DB.Main.Model(&product).Association("Taxonomies").Append(&taxonomy); err != nil {
		return err
	}
	return nil
}

func (_i *repo) DeleteAttribute(productID uint64, attrID uint64) (err error) {
	var product schema.Product
	var taxonomy schema.Taxonomy

	if err = _i.DB.Main.First(&product, productID).Error; err != nil {
		return err
	}
	if err = _i.DB.Main.First(&taxonomy, attrID).Error; err != nil {
		return err
	}

	// Remove taxonomy from the product
	if err = _i.DB.Main.Model(&product).Association("Taxonomies").Delete(&taxonomy); err != nil {
		return err
	}
	return nil
}

func (_i *repo) Delete(businessID uint64, id uint64) error {
	return _i.DB.Main.
		Where(&schema.Post{BusinessID: businessID}).
		Delete(&schema.Post{}, id).Error
}
