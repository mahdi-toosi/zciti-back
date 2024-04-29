package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Order struct {
	ID            uint64             `gorm:"primaryKey" faker:"-"`
	Status        OrderStatus        `gorm:"varchar(20); not null" faker:"oneof: pending, processing, onHold, completed, cancelled, refunded, failed"`
	TotalAmt      float64            `gorm:"not null"`
	PaymentMethod OrderPaymentMethod `gorm:"varchar(20); not null" faker:"oneof: online"`
	Meta          OrderMeta          `gorm:"type:jsonb"`
	UserID        uint64             `gorm:"index" faker:"-"`
	User          User               `gorm:"foreignKey:UserID" faker:"-"`
	BusinessID    uint64             `gorm:"index" faker:"-"`
	Business      Business           `gorm:"foreignKey:BusinessID" faker:"-"`
	ParentID      uint64             `gorm:"index" faker:"-"`
	//OrderItems    []OrderItem   `faker:"-"`
	//Transactions  []Transaction `faker:"-"`
	Base
}

type OrderStatus string

const (
	onHold     OrderStatus = "onHold"     // The order is awaiting payment or stock availability.
	failed     OrderStatus = "failed"     // The payment for the order has failed or been declined.
	pending    OrderStatus = "pending"    // The order has been placed, but no payment has been made yet.
	refunded   OrderStatus = "refunded"   // The order has been refunded to the customer.
	completed  OrderStatus = "completed"  // The order has been paid for and fulfilled.
	cancelled  OrderStatus = "cancelled"  // The order has been cancelled by the customer or the administrator.
	processing OrderStatus = "processing" // Payment has been received, and the order is being processed.
)

type OrderPaymentMethod string

const (
	PaymentMethodCash           OrderPaymentMethod = "cash"
	PaymentMethodOnline         OrderPaymentMethod = "online"
	PaymentMethodCashOnDelivery OrderPaymentMethod = "cashOnDelivery"
)

type OrderMeta struct {
	TaxAmt    uint64
	UserNote  string
	UserAgent string
	UserIP    string
}

func (bm *OrderMeta) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal OrderMeta with value %v", value)
	}
	return json.Unmarshal(byteValue, bm)
}

func (bm OrderMeta) Value() (driver.Value, error) {
	return json.Marshal(bm)
}
