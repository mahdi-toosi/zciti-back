package request

import (
	"errors"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"strings"
	"time"
)

type Coupon struct {
	ID          uint64
	Code        string            `example:"xoxoxo" validate:"required,min=1,max=100"`
	Title       string            `example:"title" validate:"required,min=1,max=255"`
	Description *string           `example:"description" validate:"omitempty,min=1,max=500"`
	Value       float64           `example:"999" validate:"required,number"`
	Type        schema.CouponType `example:"fixedAmount" validate:"required,oneof=fixedAmount percentage"`
	StartTime   string            `example:"2023-10-20T15:47:33.084Z" validate:"datetime=2006-01-02 15:04:05"`
	EndTime     string            `example:"2023-10-20T15:47:33.084Z" validate:"datetime=2006-01-02 15:04:05"`
	BusinessID  uint64            `example:"1" validate:"required"`
	Meta        schema.CouponMeta `validate:"omitempty"`
}

type Coupons struct {
	BusinessID uint64
	Title      string
	Pagination *paginator.Pagination
}

type ValidateCoupon struct {
	Code          string
	UserID        uint64
	BusinessID    uint64
	OrderTotalAmt float64
}

func (req *Coupon) ToDomain() (item *schema.Coupon, err error) {
	item = &schema.Coupon{
		ID:          req.ID,
		Type:        req.Type,
		Meta:        req.Meta,
		Title:       req.Title,
		Value:       req.Value,
		BusinessID:  req.BusinessID,
		Description: req.Description,
		Code:        strings.TrimSpace(req.Code),
	}

	if req.StartTime != "" {
		loc, _ := time.LoadLocation("Asia/Tehran")
		item.StartTime, _ = time.ParseInLocation(time.DateTime, req.StartTime, loc)
	}
	if req.EndTime != "" {
		loc, _ := time.LoadLocation("Asia/Tehran")
		item.EndTime, _ = time.ParseInLocation(time.DateTime, req.EndTime, loc)
	}

	if item.EndTime.Before(item.StartTime) {
		return nil, errors.New("تاریخ شروع پس از پایان است")
	}

	return item, nil
}
