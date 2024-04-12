package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"go-fiber-starter/utils/helpers"
	"golang.org/x/exp/slices"
)

type User struct {
	ID              uint64             `gorm:"primaryKey" faker:"-"`
	FirstName       string             `gorm:"varchar(250);" faker:"first_name"`
	LastName        string             `gorm:"varchar(250);" faker:"last_name"`
	Mobile          uint64             `gorm:"not null;uniqueIndex"`
	MobileConfirmed bool               `gorm:"default:false"`
	Permissions     UserPermissionsMap `gorm:"type:json;not null"`
	Password        string             `gorm:"varchar(250);not null"`
	Businesses      []*Business        `gorm:"many2many:business_users;" faker:"-"`
	Base
}

type UserPermissionsMap map[uint64] /* businessID*/ []UserRole

func (rm *UserPermissionsMap) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal UserPermissionsMap with value %v", value)
	}
	return json.Unmarshal(byteValue, rm)
}

func (rm UserPermissionsMap) Value() (driver.Value, error) {
	return json.Marshal(rm)
}

func (u *User) ComparePassword(password string) bool {
	return helpers.ValidateHash(password, u.Password)
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

type UserRole string

const (
	URUser          UserRole = "user"
	URAdmin         UserRole = "admin"
	URBusinessOwner UserRole = "businessOwner"
)

func (u *User) IsAdmin() bool {
	return slices.Contains(u.Permissions[ROOT_BUSINESS_ID], URAdmin)
}

func (u *User) IsBusinessOwner(businessID uint64) bool {
	roles := u.Permissions[businessID]
	return u.IsAdmin() || slices.Contains(roles, URBusinessOwner)
}
