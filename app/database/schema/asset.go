package schema

import (
	"github.com/google/uuid"
)

type Asset struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title      string    `gorm:"varchar(255);"`
	Path       string
	Ext        string
	IsPrivate  bool
	Size       uint64
	AltText    string `gorm:"varchar(255);"`
	UserID     uint64
	User       User `gorm:"foreignKey:UserID"`
	BusinessID uint64
	Business   Business `gorm:"foreignKey:BusinessID"`
	Base
}
