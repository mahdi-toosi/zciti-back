package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	Meta       ReservationMeta   `gorm:"type:jsonb"`
	Base
}

type ReservationStatus string

const (
	ReservationStatusCanceled ReservationStatus = "canceled"
	ReservationStatusReserved ReservationStatus = "reserved"
)

type UniWashCommand string

const (
	UniWashCommandON        UniWashCommand = "ON"
	UniWashCommandOFF       UniWashCommand = "OFF"
	UniWashCommandMoreWater UniWashCommand = "MORE_WATER"
	UniWashCommandOffline   UniWashCommand = "OFFLINE"
)

type ReservationMeta struct {
	UniWashLastCommand            UniWashCommand `json:",omitempty"`
	UniWashLastCommandTime        time.Time      `json:",omitempty"`
	UniWashLastCommandReferenceID string         `json:",omitempty"`
}

func (pm *ReservationMeta) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal ReservationMeta with value %v", value)
	}
	return json.Unmarshal(byteValue, pm)
}

func (pm ReservationMeta) Value() (driver.Value, error) {
	return json.Marshal(pm)
}
