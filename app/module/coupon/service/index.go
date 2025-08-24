package service

import (
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
	"github.com/gofiber/fiber/v2"
	ptime "github.com/yaa110/go-persian-calendar"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/coupon/repository"
	"go-fiber-starter/app/module/coupon/request"
	"go-fiber-starter/app/module/coupon/response"
	userRequest "go-fiber-starter/app/module/user/request"
	userService "go-fiber-starter/app/module/user/service"
	"go-fiber-starter/internal"
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

	CouponMessageSend(req request.CouponMessageSend) error
	ValidateCoupon(req request.ValidateCoupon) (coupon *schema.Coupon, err error)
	ApplyCoupon(coupon *schema.Coupon, userID uint64, totalAmt *float64) (err error)
	CalcTotalAmtWithDiscount(coupon *schema.Coupon, totalAmt *float64) (_totalAmt float64)
}

func Service(Repo repository.IRepository, userService userService.IService, messageWay *internal.MessageWayService) IService {
	return &service{
		Repo,
		userService,
		messageWay,
	}
}

type service struct {
	Repo        repository.IRepository
	UserService userService.IService
	MessageWay  *internal.MessageWayService
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

func (_i *service) CouponMessageSend(req request.CouponMessageSend) error {
	payload := userRequest.BusinessUsers{UserIDs: req.UserIDs, BusinessID: req.BusinessID}

	users, _, err := _i.UserService.Users(payload)
	if err != nil {
		return err
	}

	coupon, err := _i.Show(req.BusinessID, req.CouponID)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Tehran")
	gTime, err := time.ParseInLocation(time.DateTime, coupon.EndTime, loc)
	if err != nil {
		return err
	}
	jTime := ptime.New(gTime)
	for _, user := range users {

		_, err := _i.MessageWay.Send(MessageWay.Message{
			Provider:   5, // با سرشماره 5000
			TemplateID: 12109,
			Method:     "sms",
			Params:     []string{user.FullName, coupon.Code, jTime.Format("yyyy/MM/dd HH:mm")},
			Mobile:     fmt.Sprintf("0%d", user.Mobile),
		})
		if err != nil {
			return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "ارسال دستور با خطا مواجه شد، دوباره امتحان کنید."}
		}
	}

	return nil
}

func (_i *service) ValidateCoupon(req request.ValidateCoupon) (coupon *schema.Coupon, err error) {
	coupon, err = _i.Repo.GetOne(req.BusinessID, nil, &req.Code)
	if err != nil {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف معتبر نمی باشد"}
	}

	if slices.Contains(coupon.Meta.UsedBy, req.UserID) {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف برای شما قبلا استفاده شده است"}
	}

	now := time.Now()
	if now.After(coupon.EndTime) {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف منقضی شده است"}
	} else if now.Before(coupon.StartTime) {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "کد تخفیف در حال حاضر فعال نمی باشد"}
	}

	if coupon.Meta.LimitInReservationTime {
		jCouponStartTime := ptime.New(coupon.StartTime).Format("HH:mm - MM/dd")
		jCouponEndTime := ptime.New(coupon.EndTime).Format("HH:mm - MM/dd")

		if len(req.OrderReservationRanges) == 0 || len(req.OrderReservationRanges[0]) == 0 {
			return nil, &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: fmt.Sprintf("کد تخیف در بازه زمانی %s - %s قابل استفاده است", jCouponStartTime, jCouponEndTime),
			}
		}

		loc, _ := time.LoadLocation("Asia/Tehran")
		for _, orderReservationRang := range req.OrderReservationRanges {
			reqStartTime, err := time.ParseInLocation(time.DateTime, orderReservationRang[0], loc)
			if err != nil {
				return nil, err
			}
			reqEndTime, err := time.ParseInLocation(time.DateTime, orderReservationRang[1], loc)
			if err != nil {
				return nil, err
			}

			if reqStartTime.Before(coupon.StartTime) || reqEndTime.After(coupon.EndTime) {
				return nil, &fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: fmt.Sprintf("کد تخیف در بازه زمانی %s - %s قابل استفاده است", jCouponStartTime, jCouponEndTime),
				}
			}
		}
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
