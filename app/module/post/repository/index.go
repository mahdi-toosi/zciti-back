package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.PostsRequest) (posts []*schema.Post, paging paginator.Pagination, err error)
	GetOne(businessID uint64, id uint64) (post *schema.Post, err error)
	Create(post *schema.Post) (err error)
	Update(id uint64, post *schema.Post) (err error)
	Delete(businessID uint64, id uint64) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.PostsRequest) (posts []*schema.Post, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Post{}).Where("business_id = ?", req.BusinessID)

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Preload("Author").Preload("Business").Order("created_at desc").Find(&posts).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (post *schema.Post, err error) {
	err = _i.DB.Main.Preload("Business").
		Where("business_id = ?", businessID).
		First(&post, id).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (_i *repo) Create(post *schema.Post) (err error) {
	return _i.DB.Main.Create(post).Error
}

func (_i *repo) Update(id uint64, post *schema.Post) (err error) {
	return _i.DB.Main.Model(&schema.Post{}).
		Where(&schema.Post{ID: id}).
		Updates(post).Error
}

func (_i *repo) Delete(businessID uint64, id uint64) error {
	return _i.DB.Main.Delete(&schema.Post{}, id).Error
}
