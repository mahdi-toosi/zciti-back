package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Business struct {
	ID          uint64          `gorm:"primaryKey" faker:"-"`
	Title       string          `gorm:"not null;varchar(255)" faker:"word"`
	Type        BusinessType    `gorm:"not null;varchar(255)" faker:"oneof:GymManager,Bakery"`
	OwnerID     uint64          `gorm:"not null" faker:"-"`
	Owner       User            `gorm:"foreignKey:OwnerID" faker:"-"`
	Account     BusinessAccount `gorm:"varchar(100);default:default" faker:"-"`
	Meta        BusinessMeta    `gorm:"type:json" faker:"-"`
	Description string          `gorm:"varchar(500)" faker:"paragraph"`
	Users       []*User         `gorm:"many2many:business_users;" faker:"-"`
	Base
}

type BusinessType string

const (
	BTypeROOT          BusinessType = "ROOT"
	BTypeBakery        BusinessType = "Bakery"
	BTypeGymManager    BusinessType = "GymManager"
	BTypeWMReservation BusinessType = "WMReservation" // Washing Machine Reservation
)

var TypeDisplayProxy = map[BusinessType]string{
	BTypeGymManager: "مدیر باشگاه",
	BTypeBakery:     "نانوایی",
}

type BusinessAccount string

const (
	BusinessAccountDefault BusinessAccount = "default"
)

const ROOT_BUSINESS_ID = 1

type BusinessMeta struct {
	ShebaNumber string
	AssetsSize  uint64
}

func (bm BusinessMeta) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal BusinessMeta with value %v", value)
	}
	return json.Unmarshal(byteValue, &bm)
}

func (bm BusinessMeta) Value() (driver.Value, error) {
	return json.Marshal(bm)
}
