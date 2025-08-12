package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Coupon struct {
	ID          uint64     `gorm:"primaryKey"`
	Code        string     `gorm:"varchar(255); not null;index:idx_code,unique"`
	Title       string     `gorm:"varchar(255); not null;"`
	Description *string    `gorm:"varchar(500;"`
	Value       float64    `gorm:"not null"`
	Type        CouponType `gorm:"not null"`
	StartTime   time.Time  `gorm:"not null"`
	EndTime     time.Time  `gorm:"not null"`
	TimesUsed   int        ``
	BusinessID  uint64     `gorm:"index:idx_code; not null"`
	Business    Business   `gorm:"foreignKey:BusinessID"`
	Meta        CouponMeta `gorm:"type:jsonb"`
	Base
}

type CouponType string

const (
	CouponTypePercentage  CouponType = "percentage"
	CouponTypeFixedAmount CouponType = "fixedAmount"
)

type CouponMeta struct {
	UsedBy                 []uint64 `json:",omitempty"`
	MaxUsage               int      `json:",omitempty" validate:"required,min=1"`
	MinPrice               float64  `json:",omitempty"`
	MaxPrice               float64  `json:",omitempty"`
	MaxDiscount            float64  `json:",omitempty"`
	IncludeUserIDs         []uint64 `json:",omitempty"`
	LimitInReservationTime bool     `json:",omitempty"`
	//IncludeProducts []uint64 `json:",omitempty"`
	//ExcludeProducts []uint64 `json:",omitempty"`
}

func (cm *CouponMeta) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal CouponMeta with value %v", value)
	}
	return json.Unmarshal(byteValue, cm)
}

func (cm CouponMeta) Value() (driver.Value, error) {
	return json.Marshal(cm)
}
