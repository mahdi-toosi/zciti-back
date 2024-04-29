package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Order struct {
	ID            uint64
	Status        schema.OrderStatus        `example:"pending" validate:"omitempty,oneOf=pending processing onHold completed cancelled refunded failed"`
	PaymentMethod schema.OrderPaymentMethod `example:"page" validate:"required,oneof=product post page"`
	UserNote      string                    `example:"note note" validate:"omitempty,min=2,max=255" json:",omitempty" faker:""`
	BusinessID    uint64                    `example:"1" validate:"min=1"`
	UserID        uint64                    `example:"1" validate:"min=1"`
}

type Orders struct {
	BusinessID uint64
	Pagination *paginator.Pagination
}

func (req *Order) ToDomain() *schema.Order {
	return &schema.Order{
		ID:            req.ID,
		Status:        req.Status,
		UserID:        req.UserID,
		BusinessID:    req.BusinessID,
		PaymentMethod: req.PaymentMethod,
		Meta: schema.OrderMeta{
			UserNote: req.UserNote,
		},
	}
}
