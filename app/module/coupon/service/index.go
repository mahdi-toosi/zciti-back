package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/coupon/repository"
	"go-fiber-starter/app/module/coupon/request"
	"go-fiber-starter/app/module/coupon/response"
	"go-fiber-starter/utils/paginator"
	"golang.org/x/exp/slices"
	"time"
)

type IService interface {
	Index(req request.Coupons) (coupons []*response.Coupon, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (coupon *response.Coupon, err error)
	Store(req request.Coupon) (err error)
	Update(id uint64, req request.Coupon) (err error)
	Destroy(id uint64) error

	ValidateCoupon(req request.ValidateCoupon) (coupon *schema.Coupon, err error)
	ApplyCoupon(coupon *schema.Coupon, userID uint64, totalAmt *float64) (err error)
	CalcTotalAmtWithDiscount(coupon *schema.Coupon, totalAmt *float64) (_totalAmt float64)
}

func Service(Repo repository.IRepository) IService {
	return &service{Repo}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Index(req request.Coupons) (coupons []*response.Coupon, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		coupons = append(coupons, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(businessID uint64, id uint64) (article *response.Coupon, err error) {
	result, err := _i.Repo.GetOne(businessID, &id, nil)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req request.Coupon) (err error) {
	item, err := req.ToDomain()
	if err != nil {
		return err
	}

	return _i.Repo.Create(item)
}

func (_i *service) Update(id uint64, req request.Coupon) (err error) {
	item, err := req.ToDomain()
	if err != nil {
		return err
	}

	return _i.Repo.Update(id, item)
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}

func (_i *service) ValidateCoupon(req request.ValidateCoupon) (coupon *schema.Coupon, err error) {
	coupon, err = _i.Repo.GetOne(req.BusinessID, nil, &req.Code)
	if err != nil {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف معتبر نمی باشد"}
	}

	if slices.Contains(coupon.Meta.UsedBy, req.UserID) {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف برای شما قبلا استفاده شده است"}
	}

	loc, _ := time.LoadLocation("Asia/Tehran")
	now := time.Now().In(loc)
	if now.After(coupon.EndTime) {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف منقضی شده است"}
	} else if now.Before(coupon.StartTime) {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف در حال حاضر فعال نمی باشد"}
	}

	if coupon.Meta.MaxUsage > 0 && coupon.TimesUsed >= coupon.Meta.MaxUsage {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "تعداد استفاده از کد تخفیف بیش از حد مجاز است"}
	}

	if len(coupon.Meta.IncludeUserIDs) > 0 && !slices.Contains(coupon.Meta.IncludeUserIDs, req.UserID) {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف برای شما فعال نمی باشد"}
	}

	if coupon.Meta.MinPrice > 0 && req.OrderTotalAmt < coupon.Meta.MinPrice {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: fmt.Sprintf("مبلغ سفارش باید بیشتر از %.2f تومان باشد", coupon.Meta.MinPrice)}
	}

	if coupon.Meta.MaxPrice > 0 && req.OrderTotalAmt > coupon.Meta.MaxPrice {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: fmt.Sprintf("مبلغ سفارش باید کمتر از %.2f تومان باشد", coupon.Meta.MaxPrice)}
	}

	return coupon, nil
}

func (_i *service) CalcTotalAmtWithDiscount(coupon *schema.Coupon, totalAmt *float64) (_totalAmt float64) {
	if coupon.Type == schema.CouponTypePercentage {
		discount := *totalAmt * (coupon.Value / 100)
		if coupon.Meta.MaxDiscount != 0 && discount > coupon.Meta.MaxDiscount {
			discount = coupon.Meta.MaxDiscount
		}
		_totalAmt = *totalAmt - discount
	} else {
		_totalAmt = *totalAmt - coupon.Value
	}

	if _totalAmt < 0 {
		_totalAmt = 0
	}

	return _totalAmt
}

func (_i *service) ApplyCoupon(coupon *schema.Coupon, userID uint64, totalAmt *float64) (err error) {
	*totalAmt = _i.CalcTotalAmtWithDiscount(coupon, totalAmt)

	coupon.TimesUsed++
	coupon.Meta.UsedBy = append(coupon.Meta.UsedBy, userID)

	if err = _i.Repo.Update(coupon.ID, coupon); err != nil {
		return err
	}

	return
}
