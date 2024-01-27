package schema

import (
	"github.com/lib/pq"
	"go-fiber-starter/utils/helpers"
	"golang.org/x/exp/slices"
)

type User struct {
	ID              uint64         `gorm:"primaryKey" faker:"-"`
	FirstName       string         `gorm:"varchar(250);" faker:"first_name"`
	LastName        string         `gorm:"varchar(250);" faker:"last_name"`
	Mobile          uint64         `gorm:"not null;uniqueIndex"`
	MobileConfirmed bool           `gorm:"default:false"`
	Roles           pq.StringArray `gorm:"type:text[]"` // admin,
	Password        string         `gorm:"varchar(250);not null"`
	Businesses      []*Business    `gorm:"many2many:business_users;" faker:"-"`
	Base
}

func (u *User) ComparePassword(password string) bool {
	return helpers.ValidateHash(password, u.Password)
}

const (
	RUser          = "user"
	RAdmin         = "admin"
	RBusinessOwner = "businessOwner"
)

func (u *User) IsAdmin() bool {
	return slices.Contains(u.Roles, RAdmin)
}

func (u *User) IsBusinessOwner() bool {
	return slices.Contains(u.Roles, RBusinessOwner)
}
