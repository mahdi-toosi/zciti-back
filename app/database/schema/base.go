package schema

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	CreatedAt time.Time      `gorm:"autoCreateTime" json:",omitempty"`
	UpdatedAt time.Time      `json:",omitempty"`
	DeletedAt gorm.DeletedAt `json:",omitempty" faker:"-"`
}
