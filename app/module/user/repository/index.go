package repository

import (
	"fmt"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Users) (users []*schema.User, paging paginator.Pagination, err error)
	GetOne(id uint64) (user *schema.User, err error)
	Create(user *schema.User) (err error)
	Update(id uint64, user *schema.User) (err error)
	Delete(id uint64) (err error)
	FindUserByMobile(mobile uint64) (user *schema.User, err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Users) (users []*schema.User, paging paginator.Pagination, err error) {
	query := _i.DB.DB.Model(&schema.User{})

	if req.Keyword != "" {
		query.Where("first_name Like ?", fmt.Sprint("%", req.Keyword, "%"))
		query.Or("last_name Like ?", fmt.Sprint("%", req.Keyword, "%"))
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total
		query.Count(&total)

		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Debug().Order("created_at asc").Find(&users).Error
	err = query.Offset(req.Pagination.Offset).Limit(req.Pagination.Limit).Find(&users).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (user *schema.User, err error) {
	if err := _i.DB.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (_i *repo) Create(user *schema.User) (err error) {
	return _i.DB.DB.Create(user).Error
}

func (_i *repo) Update(id uint64, user *schema.User) (err error) {
	return _i.DB.DB.Model(&schema.User{}).
		Where(&schema.User{ID: id}).
		Updates(user).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.DB.Delete(&schema.User{}, id).Error
}

func (_i *repo) FindUserByMobile(mobile uint64) (user *schema.User, err error) {
	if err := _i.DB.DB.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
