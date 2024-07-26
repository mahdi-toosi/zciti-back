package response

import (
	"go-fiber-starter/app/database/schema"
	"time"
)

type Coupon struct {
	ID          uint64            `json:",omitempty"`
	Code        string            `json:",omitempty"`
	Title       string            `json:",omitempty"`
	Description *string           `json:",omitempty"`
	Value       float64           `json:",omitempty"`
	StartTime   string            `json:",omitempty"`
	EndTime     string            `json:",omitempty"`
	TimesUsed   int               `json:",omitempty"`
	Type        schema.CouponType `json:",omitempty"`
	Meta        schema.CouponMeta `json:",omitempty"`
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
