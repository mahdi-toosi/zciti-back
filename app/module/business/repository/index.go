package repository

import (
	"fmt"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/business/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Businesses) (businesses []*schema.Business, paging paginator.Pagination, err error)
	GetUsers(req request.Users) (users []*schema.User, paging paginator.Pagination, err error)
	InsertUser(businessID uint64, userID uint64) (err error)
	DeleteUser(businessID uint64, userID uint64) (err error)
	GetOne(id uint64) (business *schema.Business, err error)
	Create(business *schema.Business) (err error)
	Update(id uint64, business *schema.Business) (err error)
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

func (_i *repo) GetAll(req request.Businesses) (businesses []*schema.Business, paging paginator.Pagination, err error) {
	query := _i.DB.DB.Model(&schema.Business{})

	if req.Keyword != "" {
		query.Where("title Like ?", fmt.Sprint("%", req.Keyword, "%"))
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Preload("Owner").Order("created_at asc").Find(&businesses).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetUsers(req request.Users) (users []*schema.User, paging paginator.Pagination, err error) {
	query := _i.DB.DB.
		Model(&users).
		Joins("JOIN business_users ON business_users.user_id = users.id").
		Where("business_users.business_id = ?", req.BusinessID).
		Order("created_at ASC")

	if err != nil {
		return
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Find(&users).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) InsertUser(businessID uint64, userID uint64) (err error) {
	err = _i.DB.DB.
		Exec(`INSERT INTO business_users 
    		(user_id, business_id) VALUES (?, ?)`,
			userID, businessID).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (_i *repo) DeleteUser(businessID uint64, userID uint64) (err error) {
	err = _i.DB.DB.
		Exec(`DELETE FROM business_users 
       		WHERE user_id = ? AND business_id = ?`,
			userID, businessID).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (_i *repo) GetOne(id uint64) (business *schema.Business, err error) {
	if err := _i.DB.DB.First(&business, id).Error; err != nil {
		return nil, err
	}

	return business, nil
}

func (_i *repo) Create(business *schema.Business) (err error) {
	return _i.DB.DB.Create(business).Error
}

func (_i *repo) Update(id uint64, business *schema.Business) (err error) {
	return _i.DB.DB.Model(&schema.Business{}).
		Where(&schema.Business{ID: id}).
		Updates(business).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.DB.Delete(&schema.Business{}, id).Error
}