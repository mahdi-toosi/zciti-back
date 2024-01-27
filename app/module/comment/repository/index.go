package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/comment/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(postID uint64, req request.Comments) (comments []*schema.Comment, paging paginator.Pagination, err error)
	GetOne(id uint64) (comment *schema.Comment, err error)
	Create(comment *schema.Comment) (err error)
	Update(id uint64, comment *schema.Comment) (err error)
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

func (_i *repo) GetAll(postID uint64, req request.Comments) (comments []*schema.Comment, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Comment{}).Where("post_id = ?", postID)

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Preload("Author").Order("created_at desc").Find(&comments).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (comment *schema.Comment, err error) {
	if err := _i.DB.Main.First(&comment, id).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (_i *repo) Create(comment *schema.Comment) (err error) {
	return _i.DB.Main.Create(comment).Error
}

func (_i *repo) Update(id uint64, comment *schema.Comment) (err error) {
	return _i.DB.Main.Model(&schema.Comment{}).
		Where(&schema.Comment{ID: id}).
		Updates(comment).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Comment{}, id).Error
}
