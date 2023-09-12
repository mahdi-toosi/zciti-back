package schema

import (
	"github.com/bangadam/go-fiber-starter/utils/helpers"
	"github.com/lib/pq"
)

type User struct {
	ID              uint64         `gorm:"primary_key" faker:"-"`
	FirstName       string         `faker:"first_name"`
	LastName        string         `faker:"last_name"`
	Mobile          string         `gorm:"not null;uniqueIndex" faker:"e_164_phone_number"`
	MobileConfirmed bool           `gorm:"default:false"`
	Roles           pq.StringArray `gorm:"type:text[]" faker:"slice_len=2"`
	Password        string         `gorm:"not null" faker:"password"`
	Base
}

func (u *User) ComparePassword(password string) bool {
	return helpers.ValidateHash(password, u.Password)
}
