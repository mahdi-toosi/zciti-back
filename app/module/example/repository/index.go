package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/example/request"
	"go-fiber-starter/app/module/example/response"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Examples) (examples []*response.Example, paging paginator.Pagination, err error)
	GetOne(businessID uint64, id uint64) (example *response.Example, err error)
	Create(example *schema.Example) (err error)
	Update(id uint64, example *schema.Example) (err error)
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

func (_i *repo) GetAll(req request.Examples) (examples []*response.Example, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Example{}).
		Where(&schema.Example{BusinessID: req.BusinessID})

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at desc").Find(&examples).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (example *response.Example, err error) {
	if err := _i.DB.Main.
		Where(&schema.Example{BusinessID: businessID}).
		First(&example, id).
		Error; err != nil {
		return nil, err
	}

	return example, nil
}

func (_i *repo) Create(example *schema.Example) (err error) {
	return _i.DB.Main.Create(example).Error
}

func (_i *repo) Update(id uint64, example *schema.Example) (err error) {
	return _i.DB.Main.Model(&schema.Example{}).
		Where(&schema.Example{ID: id, BusinessID: example.BusinessID}).
		Updates(example).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Example{}, id).Error
}
