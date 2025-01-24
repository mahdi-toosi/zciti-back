package repository

import (
	"errors"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/coupon/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
	"strings"
)

type IRepository interface {
	GetAll(req request.Coupons) (coupons []*schema.Coupon, paging paginator.Pagination, err error)
	GetOne(businessID uint64, id *uint64, code *string) (coupon *schema.Coupon, err error)
	Create(coupon *schema.Coupon) (err error)
	Update(id uint64, coupon *schema.Coupon) (err error)
	Delete(id uint64) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{DB}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Coupons) (coupons []*schema.Coupon, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Coupon{}).
		Where(&schema.Coupon{BusinessID: req.BusinessID})

	if req.Title != "" {
		query.Where("title LIKE ?", "%"+req.Title+"%")
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at desc").Find(&coupons).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id *uint64, code *string) (coupon *schema.Coupon, err error) {
	query := _i.DB.Main.
		Model(&schema.Coupon{}).
		Where(&schema.Coupon{BusinessID: businessID})

	if id != nil {
		if err = query.First(&coupon, id).Error; err != nil {
			return nil, err
		}
	}
	if code != nil {
		query.Where("code = ?", code)

		if err = query.First(&coupon).Error; err != nil {
			return nil, err
		}
	}

	return coupon, nil
}

func (_i *repo) Create(coupon *schema.Coupon) (err error) {
	err = _i.DB.Main.Create(coupon).Error

	if err != nil && strings.Contains(err.Error(), "value violates unique constraint") {
		return errors.New("این کد کوپن قبلا ثبت شده است، لطفا مقداری خاص ثبت کنید")
	}

	return err
}

func (_i *repo) Update(id uint64, coupon *schema.Coupon) (err error) {
	err = _i.DB.Main.Model(&schema.Coupon{}).
		Where(&schema.Coupon{ID: id, BusinessID: coupon.BusinessID}).
		Updates(coupon).Error

	if err != nil && strings.Contains(err.Error(), "value violates unique constraint") {
		return errors.New("این کد کوپن قبلا ثبت شده است، لطفا مقداری خاص ثبت کنید")
	}

	return err
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Coupon{}, id).Error
}
