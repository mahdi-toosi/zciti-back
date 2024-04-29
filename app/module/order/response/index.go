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

func FromDomain(item *schema.Order) (res *Order) {
	if item == nil {
		return nil
	}

	return &Order{
		ID:            item.ID,
		Meta:          item.Meta,
		Status:        item.Status,
		ParentID:      item.ParentID,
		TotalAmt:      item.TotalAmt,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		PaymentMethod: item.PaymentMethod,
		User: response.User{
			ID:       item.User.ID,
			FullName: item.User.FullName(),
		},
	}
}
