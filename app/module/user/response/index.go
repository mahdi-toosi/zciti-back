package response

import (
	"go-fiber-starter/app/database/schema"
)

type User struct {
	ID              uint64
	FirstName       string
	LastName        string
	Mobile          uint64
	MobileConfirmed bool
	Roles           []string
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
		}
	}

	return res
}
