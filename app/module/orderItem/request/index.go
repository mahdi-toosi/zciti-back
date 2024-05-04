package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type OrderItem struct {
	ID        uint64
	Quantity  int    `example:"1" validate:"number,min=1"`
	ProductID uint64 `example:"1" validate:"number,min=1"`
	OrderID   uint64 `example:"1" validate:"number,min=1"`
}

type OrderItems struct {
	Pagination *paginator.Pagination
}

func (req *OrderItem) ToDomain() *schema.OrderItem {
	return &schema.OrderItem{
		ID:        req.ID,
		OrderID:   req.OrderID,
		Quantity:  req.Quantity,
		ProductID: req.ProductID,
	}
}
