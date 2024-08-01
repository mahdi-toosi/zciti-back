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
	Permissions     schema.UserPermissions ``
	Roles           []schema.UserRole      `json:",omitempty"`
}

func FromDomain(item *schema.User, businessID *uint64) (res *User) {
	if item == nil {
		return nil
	}

	res = &User{
		ID:              item.ID,
		Mobile:          item.Mobile,
		LastName:        item.LastName,
		FirstName:       item.FirstName,
		FullName:        item.FullName(),
		MobileConfirmed: item.MobileConfirmed,
	}

	if businessID != nil {
		res.Roles = item.Permissions[*businessID]
	} else {
		res.Permissions = item.Permissions
	}
	return res
}
