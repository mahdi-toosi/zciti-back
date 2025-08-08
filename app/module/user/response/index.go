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
	CityID          uint64                 `json:",omitempty"`
	CityTitle       string                 `json:",omitempty"`
	WorkspaceID     uint64                 `json:",omitempty"`
	WorkspaceTitle  string                 `json:",omitempty"`
	DormitoryID     uint64                 `json:",omitempty"`
	DormitoryTitle  string                 `json:",omitempty"`
	IsSuspended     bool                   ``
	SuspenseReason  string                 `json:",omitempty"`
	Permissions     schema.UserPermissions `json:",omitempty"`
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
		IsSuspended:     *item.IsSuspended,
		MobileConfirmed: item.MobileConfirmed,
	}

	if businessID != nil {
		res.Roles = item.Permissions[*businessID]
		if item.SuspenseReason != nil {
			res.SuspenseReason = *item.SuspenseReason
		}

		if item.City != nil {
			res.CityTitle = item.City.Title
			res.CityID = item.City.ID
		}

		if item.Workspace != nil {
			res.WorkspaceTitle = item.Workspace.Title
			res.WorkspaceID = item.Workspace.ID
		}

		if item.Dormitory != nil {
			res.DormitoryTitle = item.Dormitory.Title
			res.DormitoryID = item.Dormitory.ID
		}
	} else {
		res.Permissions = item.Permissions
	}
	return res
}
