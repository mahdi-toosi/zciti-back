package schema

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
