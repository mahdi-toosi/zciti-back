package response

import (
	"go-fiber-starter/app/database/schema"
)

type User struct {
	ID              uint64                 `json:",omitempty"`
	FirstName       string                 `json:",omitempty"`
	LastName        string                 `json:",omitempty"`
	FullName        string                 `json:",omitempty"`
	Mobile          uint64                 `json:",omitempty"`
	MobileConfirmed bool                   `json:",omitempty"`
	Permissions     schema.UserPermissions `json:",omitempty"`
}

func FromDomain(item *schema.User) (res *User) {
	if item != nil {
		res = &User{
			ID:              item.ID,
			FullName:        item.FullName(),
			Mobile:          item.Mobile,
			MobileConfirmed: item.MobileConfirmed,
			Permissions:     item.Permissions,
		}
	}

	return res
}
