package response

import (
	"go-fiber-starter/app/database/schema"
)

type User struct {
	ID               uint64                 `json:",omitempty"`
	FirstName        string                 `json:",omitempty"`
	LastName         string                 `json:",omitempty"`
	FullName         string                 `json:",omitempty"`
	Mobile           uint64                 `json:",omitempty"`
	MobileConfirmed  bool                   `json:",omitempty"`
	CityID           uint64                 `json:",omitempty"`
	CityTitle        string                 `json:",omitempty"`
	WorkspaceID      uint64                 `json:",omitempty"`
	WorkspaceTitle   string                 `json:",omitempty"`
	DormitoryID      uint64                 `json:",omitempty"`
	DormitoryTitle   string                 `json:",omitempty"`
	ReservationCount uint64                 `json:",omitempty"`
	IsSuspended      bool                   ``
	SuspenseReason   string                 `json:",omitempty"`
	Permissions      schema.UserPermissions `json:",omitempty"`
	Roles            []schema.UserRole      `json:",omitempty"`
}

func FromDomain(item *schema.User, businessID *uint64) (res *User) {
	if item == nil {
		return nil
	}

	res = &User{
		ID:               item.ID,
		Mobile:           item.Mobile,
		LastName:         item.LastName,
		FirstName:        item.FirstName,
		FullName:         item.FullName(),
		IsSuspended:      *item.IsSuspended,
		MobileConfirmed:  item.MobileConfirmed,
		ReservationCount: item.ReservationCount,
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

type BusinessUser struct {
	ID               uint64                 `gorm:"primaryKey" faker:"-"`
	FirstName        string                 `gorm:"varchar(255);" faker:"first_name"`
	LastName         string                 `gorm:"varchar(255);" faker:"last_name"`
	Mobile           uint64                 `gorm:"not null;uniqueIndex"`
	MobileConfirmed  bool                   `gorm:"default:false"`
	ShowMobile       bool                   ``
	IsSuspended      *bool                  `gorm:"default:false"`
	SuspenseReason   *string                `gorm:"varchar(500);"`
	Permissions      schema.UserPermissions `gorm:"type:jsonb;not null"`
	Password         string                 `gorm:"varchar(255);not null"`
	CityID           *uint64                `gorm:"" faker:"-"`
	City             *schema.Taxonomy       `gorm:"foreignKey:CityID" faker:"-"`
	WorkspaceID      *uint64                `gorm:"" faker:"-"`
	Workspace        *schema.Taxonomy       `gorm:"foreignKey:WorkspaceID" faker:"-"`
	DormitoryID      *uint64                `gorm:"" faker:"-"`
	Dormitory        *schema.Taxonomy       `gorm:"foreignKey:DormitoryID" faker:"-"`
	Businesses       []*schema.Business     `gorm:"many2many:business_users;" faker:"-"`
	ReservationCount uint64
}

func FromDomainForBusinessUser(item *BusinessUser, businessID *uint64) (res *User) {
	if item == nil {
		return nil
	}

	res = &User{
		ID:               item.ID,
		Mobile:           item.Mobile,
		LastName:         item.LastName,
		FirstName:        item.FirstName,
		IsSuspended:      *item.IsSuspended,
		MobileConfirmed:  item.MobileConfirmed,
		ReservationCount: item.ReservationCount,
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
