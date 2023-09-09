package repository

import (
	"github.com/bangadam/go-fiber-starter/app/database/schema"
	"github.com/bangadam/go-fiber-starter/app/module/user/request"
	"github.com/bangadam/go-fiber-starter/internal/bootstrap/database"
	"github.com/bangadam/go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetUsers(req request.UsersRequest) (users []*schema.User, paging paginator.Pagination, err error)
	FindOne(id uint64) (user *schema.User, err error)
	Create(user *schema.User) (err error)
	Update(id uint64, user *schema.User) (err error)
	Delete(id uint64) (err error)
	FindUserByUsername(username string) (user *schema.User, err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetUsers(req request.UsersRequest) (users []*schema.User, paging paginator.Pagination, err error) {
	var count int64

	query := _i.DB.DB.Model(&schema.User{})
	query.Count(&count)

	req.Pagination.Count = count
	req.Pagination = paginator.Paging(req.Pagination)

	err = query.Offset(req.Pagination.Offset).Limit(req.Pagination.Limit).Find(&users).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) FindOne(id uint64) (article *schema.User, err error) {
	if err := _i.DB.DB.First(&article, id).Error; err != nil {
		return nil, err
	}

	return article, nil
}

func (_i *repo) Create(article *schema.User) (err error) {
	return _i.DB.DB.Create(article).Error
}

func (_i *repo) Update(id uint64, article *schema.User) (err error) {
	return _i.DB.DB.Model(&schema.User{}).
		Where(&schema.User{ID: id}).
		Updates(article).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.DB.Delete(&schema.User{}, id).Error
}

func (_i *repo) FindUserByUsername(username string) (user *schema.User, err error) {
	if err := _i.DB.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
