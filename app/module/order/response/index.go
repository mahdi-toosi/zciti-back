package response

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/response"
	"time"
)

type Order struct {
	ID            uint64
	ParentID      uint64                    `json:",omitempty"`
	TotalAmt      float64                   `json:",omitempty"`
	CreatedAt     time.Time                 `json:",omitempty"`
	UpdatedAt     time.Time                 `json:",omitempty"`
	User          response.User             `json:",omitempty"`
	Meta          schema.OrderMeta          `json:",omitempty"`
	Status        schema.OrderStatus        `json:",omitempty"`
	PaymentMethod schema.OrderPaymentMethod `json:",omitempty"`
}

func FromDomain(order *schema.Order) (res *Order) {
	if order == nil {
		return nil
	}

	return &Order{
		ID:            order.ID,
		Meta:          order.Meta,
		Status:        order.Status,
		ParentID:      order.ParentID,
		TotalAmt:      order.TotalAmt,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
		PaymentMethod: order.PaymentMethod,
		User: response.User{
			ID:       order.User.ID,
			FullName: order.User.FullName(),
		},
	}
}
