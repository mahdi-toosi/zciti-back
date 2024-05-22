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
	ParentID      *uint64            `faker:"-"`
	OrderItems    []OrderItem        `faker:"-"`
	//Transactions  []Transaction `faker:"-"`
	Base
}

type OrderStatus string

const (
	OrderStatusOnHold     OrderStatus = "onHold"     // The order is awaiting payment or stock availability.
	OrderStatusFailed     OrderStatus = "failed"     // The payment for the order has failed or been declined.
	OrderStatusPending    OrderStatus = "pending"    // The order has been placed, but no payment has been made yet.
	OrderStatusRefunded   OrderStatus = "refunded"   // The order has been refunded to the customer.
	OrderStatusCompleted  OrderStatus = "completed"  // The order has been paid for and fulfilled.
	OrderStatusCancelled  OrderStatus = "cancelled"  // The order has been cancelled by the customer or the administrator.
	OrderStatusProcessing OrderStatus = "processing" // Payment has been received, and the order is being processed.
)

type OrderPaymentMethod string

const (
	OrderPaymentMethodCash           OrderPaymentMethod = "cash"
	OrderPaymentMethodOnline         OrderPaymentMethod = "online"
	OrderPaymentMethodCashOnDelivery OrderPaymentMethod = "cashOnDelivery"
)

type OrderMeta struct {
	UserIP           string
	TaxAmt           uint64
	UserNote         string
	UserAgent        string
	PaymentAuthority string
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
