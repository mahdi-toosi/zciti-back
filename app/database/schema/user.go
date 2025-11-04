package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"go-fiber-starter/utils/helpers"
	"golang.org/x/exp/slices"
)

type User struct {
	ID               uint64          `gorm:"primaryKey" faker:"-"`
	FirstName        string          `gorm:"varchar(255);" faker:"first_name"`
	LastName         string          `gorm:"varchar(255);" faker:"last_name"`
	Mobile           uint64          `gorm:"not null;uniqueIndex"`
	MobileConfirmed  bool            `gorm:"default:false"`
	ShowMobile       bool            ``
	IsSuspended      *bool           `gorm:"default:false"`
	SuspenseReason   *string         `gorm:"varchar(500);"`
	Permissions      UserPermissions `gorm:"type:jsonb;not null"`
	Password         string          `gorm:"varchar(255);not null"`
	CityID           *uint64         `gorm:"" faker:"-"`
	City             *Taxonomy       `gorm:"foreignKey:CityID" faker:"-"`
	WorkspaceID      *uint64         `gorm:"" faker:"-"`
	Workspace        *Taxonomy       `gorm:"foreignKey:WorkspaceID" faker:"-"`
	DormitoryID      *uint64         `gorm:"" faker:"-"`
	Dormitory        *Taxonomy       `gorm:"foreignKey:DormitoryID" faker:"-"`
	Businesses       []*Business     `gorm:"many2many:business_users;" faker:"-"`
	ReservationCount uint64          `gorm:"" faker:"-"`
	Meta             *UserMeta       `gorm:"type:jsonb" faker:"-"`
	//FullName  string `gorm:"->;type:GENERATED ALWAYS AS (concat(first_name,' ',last_name));default:(-);"`
	Base
}

type UserPermissions map[uint64] /* businessID*/ []UserRole

func (up *UserPermissions) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal UserPermissions with value %v", value)
	}
	return json.Unmarshal(byteValue, up)
}

func (up UserPermissions) Value() (driver.Value, error) {
	return json.Marshal(up)
}

func (u *User) ComparePassword(password string) bool {
	return helpers.ValidateHash(password, u.Password)
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

type UserRole string

const (
	URUser             UserRole = "user"
	URAdmin            UserRole = "admin"
	URBusinessOwner    UserRole = "businessOwner"
	URBusinessObserver UserRole = "businessObserver"
)

func (u *User) IsAdmin() bool {
	return slices.Contains(u.Permissions[ROOT_BUSINESS_ID], URAdmin)
}

func (u *User) IsObserver(BusinessID uint64) bool {
	return slices.Contains(u.Permissions[BusinessID], URBusinessObserver)
}

func (u *User) IsBusinessOwner(businessID uint64) bool {
	roles := u.Permissions[businessID]
	return u.IsAdmin() || slices.Contains(roles, URBusinessOwner)
}

type UserMetaTaxonomiesToObserve map[uint64]struct {
	Checked        bool `json:"checked"`
	PartialChecked bool `json:"partialChecked"`
}

type UserMeta struct {
	PostsToObserve      []uint64                    `json:",omitempty" example:"[1,2,3]"`
	TaxonomiesToObserve UserMetaTaxonomiesToObserve `json:",omitempty" example:"{1: { checked: true; partialChecked: false }}"`
}

func (um *UserMeta) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal OrderItemMeta with value %v", value)
	}
	return json.Unmarshal(byteValue, um)
}

func (um UserMeta) Value() (driver.Value, error) {
	return json.Marshal(um)
}

func (um UserMeta) GetTaxonomiesToObserve(checked bool, partial bool) (arr []uint64) {
	if len(um.TaxonomiesToObserve) == 0 {
		return arr
	}
	for id, t := range um.TaxonomiesToObserve {
		if checked && t.Checked {
			arr = append(arr, id)
		}
		if partial && t.PartialChecked {
			arr = append(arr, id)
		}
	}
	return arr
}
