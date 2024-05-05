package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type OrderItem struct {
	ID            uint64        `gorm:"primaryKey" faker:"-"`
	Type          OrderItemType `gorm:"varchar(50); not null;" faker:"oneof: lineItem, reservation, fee, tax, coupon, shipping"`
	Quantity      int           `gorm:"not null"` // The quantity of the product ordered.
	Price         float64       `gorm:"not null"` // The price of the product at the time of the order.
	Subtotal      float64       `gorm:"not null"` // The subtotal for the order item (quantity * price).
	TaxAmt        float64       `gorm:"not null"`
	ReservationID *uint64       `faker:"-"`
	Reservation   Reservation   `gorm:"foreignKey:ReservationID" faker:"-"`
	PostID        uint64        `faker:"-"`
	Post          Post          `gorm:"foreignKey:PostID" faker:"-"`
	OrderID       uint64        `gorm:"index" faker:"-"`
	Order         Order         `gorm:"foreignKey:OrderID" faker:"-"`
	Meta          OrderItemMeta `gorm:"type:jsonb"`
	Base
}

type OrderItemType string

const (
	OrderItemTypeFee         OrderItemType = "fee"      // This order item type represents a fee that was added to the order. Fees can be added by plugins or custom code, and they can be used for a variety of purposes, such as charging for gift wrapping or rush processing.
	OrderItemTypeTax         OrderItemType = "tax"      // This order item type represents a tax charge. There is usually one tax line item per tax rate that applies to the order, but there can be more if the order contains products with different tax rates.
	OrderItemTypeCoupon      OrderItemType = "coupon"   // This order item type represents a coupon that was applied to the order. Each coupon code that is used in the order will have a corresponding coupon line item.
	OrderItemTypeLineItem    OrderItemType = "lineItem" // This is the most common order item type, and it represents a product that was ordered. Each line item corresponds to one product in the order.
	OrderItemTypeShipping    OrderItemType = "shipping" // This order item type represents a shipping charge. There is usually one shipping line item per order, but there can be more if the order is split into multiple shipments.
	OrderItemTypeReservation OrderItemType = "reservation"
)

type OrderItemMeta struct {
	TaxAmt        uint64
	ProductID     uint64
	ProductTitle  string
	ReservationID uint64
}

func (oim *OrderItemMeta) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal OrderItemMeta with value %v", value)
	}
	return json.Unmarshal(byteValue, oim)
}

func (oim OrderItemMeta) Value() (driver.Value, error) {
	return json.Marshal(oim)
}
