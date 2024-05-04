package response

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/product/response"
)

type OrderItem struct {
	ID       uint64
	Quantity int                  `json:",omitempty"`
	Price    float64              `json:",omitempty"`
	Subtotal float64              `json:",omitempty"`
	TaxAmt   float64              `json:",omitempty"`
	Product  response.Product     `json:",omitempty"`
	Type     schema.OrderItemType `json:",omitempty"`
	Meta     schema.OrderItemMeta `json:",omitempty"`
}

func FromDomain(item *schema.OrderItem) (res *OrderItem) {
	if item == nil {
		return nil
	}

	return &OrderItem{
		ID:       item.ID,
		Type:     item.Type,
		Meta:     item.Meta,
		Price:    item.Price,
		TaxAmt:   item.TaxAmt,
		Quantity: item.Quantity,
		Subtotal: item.Subtotal,
		//Product:  item.Product,
	}
}
