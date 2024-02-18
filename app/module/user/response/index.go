package response

import (
	"go-fiber-starter/app/database/schema"
)

type User struct {
	ID              uint64   `json:",omitempty"`
	FirstName       string   `json:",omitempty"`
	LastName        string   `json:",omitempty"`
	FullName        string   `json:",omitempty"`
	Mobile          uint64   `json:",omitempty"`
	MobileConfirmed bool     `json:",omitempty"`
	Roles           []string `json:",omitempty"`
}

func FromDomain(user *schema.User) (res *User) {
	if user != nil {
		res = &User{
			ID:              user.ID,
			FullName:        user.FirstName + " " + user.LastName,
			Mobile:          user.Mobile,
			MobileConfirmed: user.MobileConfirmed,
			Roles:           user.Roles,
		}
	}

	return res
}
