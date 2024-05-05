package schema

import (
	"time"
)

// Reservation TODO add proper index
type Reservation struct {
	ID         uint64            `gorm:"primaryKey" faker:"-"`
	Status     ReservationStatus `gorm:"default:reserved"`
	StartTime  time.Time         `gorm:"not null" faker:"-"`
	EndTime    time.Time         `gorm:"not null" faker:"-"`
	UserID     uint64            `gorm:"not null" faker:"-"`
	User       User              `gorm:"foreignKey:UserID" faker:"-"`
	ProductID  uint64            `gorm:"not null;index" faker:"-"`
	Product    Product           `gorm:"foreignKey:ProductID" faker:"-"`
	BusinessID uint64            `gorm:"not null" faker:"-"`
	Business   Business          `gorm:"foreignKey:BusinessID" faker:"-"`
	Base
}

type ReservationStatus string

const (
	ReservationStatusCanceled ReservationStatus = "canceled"
	ReservationStatusReserved ReservationStatus = "reserved"
)
