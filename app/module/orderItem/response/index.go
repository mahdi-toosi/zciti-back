package response

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/uniwash/response"
)

type OrderItem struct {
	ID          uint64
	Quantity    int                   `json:",omitempty"`
	PostID      uint64                `json:",omitempty"`
	Price       float64               `json:",omitempty"`
	Subtotal    float64               `json:",omitempty"`
	TaxAmt      float64               `json:",omitempty"`
	Reservation *response.Reservation `json:",omitempty"`
	Type        schema.OrderItemType  `json:",omitempty"`
	Meta        schema.OrderItemMeta  `json:",omitempty"`
}

func FromDomain(item *schema.OrderItem) (res *OrderItem) {
	if item == nil {
		return nil
	}

	oi := &OrderItem{
		ID:       item.ID,
		Type:     item.Type,
		Meta:     item.Meta,
		Price:    item.Price,
		PostID:   item.PostID,
		TaxAmt:   item.TaxAmt,
		Quantity: item.Quantity,
		Subtotal: item.Subtotal,
		//Product:  item.Product,
	}
	if item.Reservation != nil {
		oi.Reservation = &response.Reservation{
			ID:          item.Reservation.ID,
			EndTime:     item.Reservation.EndTime,
			StartTime:   item.Reservation.StartTime,
			ProductID:   item.Reservation.ProductID,
			LastCommand: item.Reservation.Meta.UniWashLastCommand,
			//ProductSKU:    item.Reservation.ProductSKU,
			//ProductTitle:  item.Reservation.ProductTitle,
			//ProductDetail: item.Reservation.ProductDetail,
		}
	}
	return oi
}
