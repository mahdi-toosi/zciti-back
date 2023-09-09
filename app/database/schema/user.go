package schema

import (
	"github.com/bangadam/go-fiber-starter/utils/helpers"
)

type User struct {
	ID              uint64    `gorm:"primary_key;column:id" json:"id"`
	FirstName       *string   `gorm:"column:first_name;default:null" json:"firstName"`
	LastName        *string   `gorm:"column:last_name;default:null" json:"lastName"`
	Mobile          string    `gorm:"column:mobile" json:"mobile"`
	MobileConfirmed *bool     `gorm:"column:mobile_confirmed;default:false" json:"mobileConfirmed"`
	Roles           *[]string `gorm:"column:roles;type:text[]" json:"roles"`
	Password        *string   `gorm:"column:password"`
	Base
}

// compare password
func (u *User) ComparePassword(password string) bool {
	return helpers.ValidateHash(password, *u.Password)
}
