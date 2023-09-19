package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.PostsRequest) (posts []*schema.Post, paging paginator.Pagination, err error)
	GetOne(id uint64) (post *schema.Post, err error)
	Create(post *schema.Post) (err error)
	Update(id uint64, post *schema.Post) (err error)
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

func (_i *repo) GetAll(req request.PostsRequest) (posts []*schema.Post, paging paginator.Pagination, err error) {
	var count int64

	query := _i.DB.DB.Model(&schema.Post{})
	query.Count(&count)

	req.Pagination.Count = count
	req.Pagination = paginator.Paging(req.Pagination)

	err = query.Offset(req.Pagination.Offset).Limit(req.Pagination.Limit).Find(&posts).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (post *schema.Post, err error) {
	if err := _i.DB.DB.First(&post, id).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (_i *repo) Create(post *schema.Post) (err error) {
	return _i.DB.DB.Create(post).Error
}

func (_i *repo) Update(id uint64, post *schema.Post) (err error) {
	return _i.DB.DB.Model(&schema.Post{}).
		Where(&schema.Post{ID: id}).
		Updates(post).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.DB.Delete(&schema.Post{}, id).Error
}
