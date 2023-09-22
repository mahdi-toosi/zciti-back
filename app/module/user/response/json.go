package response

import (
	"time"

	"go-fiber-starter/app/database/schema"
)

type User struct {
	ID              uint64   `json:"id"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	Mobile          uint64   `json:"mobile"`
	MobileConfirmed bool     `json:"mobileConfirmed"`
	Roles           []string `json:"roles"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
