package response

import (
	"go-fiber-starter/app/database/schema"
	"time"
)

type Coupon struct {
	ID          uint64
	Code        string
	Title       string
	Description *string
	Value       float64
	Type        schema.CouponType
	StartTime   string
	EndTime     string
	TimesUsed   int
	Meta        schema.CouponMeta
}

func FromDomain(item *schema.Coupon) (res *Coupon) {
	if item == nil {
		return nil
	}
	res = &Coupon{
		ID:          item.ID,
		Type:        item.Type,
		Code:        item.Code,
		Meta:        item.Meta,
		Title:       item.Title,
		Value:       item.Value,
		TimesUsed:   item.TimesUsed,
		Description: item.Description,
	}

	res.EndTime = item.EndTime.Format(time.DateTime)
	res.StartTime = item.StartTime.Format(time.DateTime)

	return res
}
