package schema

import (
	"github.com/bangadam/go-fiber-starter/utils/helpers"
	"github.com/lib/pq"
)

type User struct {
	ID              uint64         `gorm:"primary_key;column:id" json:"id" faker:"-"`
	FirstName       string         `gorm:"column:first_name" json:"firstName" faker:"first_name"`
	LastName        string         `gorm:"column:last_name" json:"lastName" faker:"last_name"`
	Mobile          string         `gorm:"column:mobile;not null;uniqueIndex" json:"mobile" faker:"e_164_phone_number"`
	MobileConfirmed bool           `gorm:"column:mobile_confirmed;default:false" json:"mobileConfirmed"`
	Roles           pq.StringArray `gorm:"column:roles;type:text[]" json:"roles" faker:"slice_len=2"`
	Password        string         `gorm:"column:password;not null" faker:"password"`
	Base
}

func (u *User) ComparePassword(password string) bool {
	return helpers.ValidateHash(password, u.Password)
}
