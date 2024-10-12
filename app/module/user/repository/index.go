package repository

import (
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

	GetUsers(req request.BusinessUsers) (users []*schema.User, paging paginator.Pagination, err error)
	GetUser(req request.BusinessUsersStoreRole) (user *schema.User, err error)
	InsertUser(businessID uint64, userID uint64) (err error)
	DeleteUser(businessID uint64, userID uint64) (err error)
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
	query := _i.DB.Main.Model(&schema.User{})

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

	err = query.Order("created_at desc").Find(&users).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (user *schema.User, err error) {
	if err := _i.DB.Main.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (_i *repo) Create(user *schema.User) (err error) {
	return _i.DB.Main.Create(user).Error
}

func (_i *repo) Update(id uint64, user *schema.User) (err error) {
	return _i.DB.Main.Model(&schema.User{}).
		Where(&schema.User{ID: id}).
		Updates(user).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.User{}, id).Error
}

func (_i *repo) FindUserByMobile(mobile uint64) (user *schema.User, err error) {
	if err := _i.DB.Main.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (_i *repo) GetUsers(req request.BusinessUsers) (users []*schema.User, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&users).
		Joins("JOIN business_users ON business_users.user_id = users.id").
		Where("business_users.business_id = ?", req.BusinessID).
		Order("created_at ASC")

	if len(req.UserIDs) > 0 {
		query.Where("users.id IN (?)", req.UserIDs)
	}

	if req.Pagination != nil && req.Pagination.Page > 0 {
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

	if req.Pagination != nil {
		paging = *req.Pagination
	}

	return
}

func (_i *repo) GetUser(req request.BusinessUsersStoreRole) (user *schema.User, err error) {
	if err := _i.DB.Main.Exec(
		`SELECT FROM business_users WHERE user_id = ? AND business_id = ?`,
		req.UserID, req.BusinessID,
	).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (_i *repo) InsertUser(businessID uint64, userID uint64) (err error) {
	err = _i.DB.Main.
		Exec(`INSERT INTO business_users (user_id, business_id) VALUES (?, ?)`,
			userID, businessID,
		).Error

	if err != nil {
		return err
	}

	return nil
}

func (_i *repo) DeleteUser(businessID uint64, userID uint64) (err error) {
	err = _i.DB.Main.
		Exec(
			`DELETE FROM business_users WHERE user_id = ? AND business_id = ?`,
			userID, businessID,
		).Error

	if err != nil {
		return err
	}

	return nil
}
