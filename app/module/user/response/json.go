package response

import (
	"time"

	"go-fiber-starter/app/database/schema"
)

type User struct {
	ID              uint64
	FirstName       string
	LastName        string
	Mobile          uint64
	MobileConfirmed bool
	Roles           []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromDomain(user *schema.User) (res *User) {
	if user != nil {
		res = &User{
			ID:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Mobile:          user.Mobile,
			MobileConfirmed: user.MobileConfirmed,
			Roles:           user.Roles,

			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return res
}
