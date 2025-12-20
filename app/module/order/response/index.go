package response

import (
	"go-fiber-starter/app/database/schema"
	cresponse "go-fiber-starter/app/module/coupon/response"
	oresponse "go-fiber-starter/app/module/orderItem/response"
	"go-fiber-starter/app/module/user/response"
	"go-fiber-starter/utils/paginator"
	"time"
)

type Order struct {
	ID            uint64
	BusinessID    uint64                    `json:",omitempty"`
	ParentID      *uint64                   `json:",omitempty"`
	TotalAmt      float64                   `json:",omitempty"`
	CreatedAt     time.Time                 `json:",omitempty"`
	UpdatedAt     time.Time                 `json:",omitempty"`
	User          response.User             `json:",omitempty"`
	Meta          schema.OrderMeta          `json:",omitempty"`
	Coupon        cresponse.Coupon          `json:",omitempty"`
	Status        schema.OrderStatus        `json:",omitempty"`
	PaymentMethod schema.OrderPaymentMethod `json:",omitempty"`
	OrderItems    []oresponse.OrderItem     `json:",omitempty"`
}

type Orders struct {
	TotalAmount uint64
	Meta        paginator.Pagination `json:",omitempty"`
}

func FromDomain(item *schema.Order) (res *Order) {
	if item == nil {
		return nil
	}

	o := &Order{
		ID:            item.ID,
		Meta:          item.Meta,
		Status:        item.Status,
		ParentID:      item.ParentID,
		TotalAmt:      item.TotalAmt,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		BusinessID:    item.BusinessID,
		PaymentMethod: item.PaymentMethod,
		Coupon: cresponse.Coupon{
			ID:    item.Coupon.ID,
			Title: item.Coupon.Title,
		},
		User: response.User{
			ID:       item.User.ID,
			Mobile:   item.User.Mobile,
			FullName: item.User.FullName(),
		},
	}

	for _, orderItem := range item.OrderItems {
		oi := oresponse.FromDomain(&orderItem)
		o.OrderItems = append(o.OrderItems, *oi)
	}

	return o
}
