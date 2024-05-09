package response

import (
	"go-fiber-starter/app/database/schema"
)

type OrderItem struct {
	ID            uint64
	Quantity      int                  `json:",omitempty"`
	PostID        uint64               `json:",omitempty"`
	Price         float64              `json:",omitempty"`
	Subtotal      float64              `json:",omitempty"`
	TaxAmt        float64              `json:",omitempty"`
	ReservationID *uint64              `json:",omitempty"`
	Type          schema.OrderItemType `json:",omitempty"`
	Meta          schema.OrderItemMeta `json:",omitempty"`
}

func FromDomain(item *schema.OrderItem) (res *OrderItem) {
	if item == nil {
		return nil
	}

	return &OrderItem{
		ID:            item.ID,
		Type:          item.Type,
		Meta:          item.Meta,
		Price:         item.Price,
		PostID:        item.PostID,
		TaxAmt:        item.TaxAmt,
		Quantity:      item.Quantity,
		Subtotal:      item.Subtotal,
		ReservationID: item.ReservationID,
		//Product:  item.Product,
	}
}
