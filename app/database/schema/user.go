package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"go-fiber-starter/utils/helpers"
	"golang.org/x/exp/slices"
)

type User struct {
	ID              uint64          `gorm:"primaryKey" faker:"-"`
	FirstName       string          `gorm:"varchar(255);" faker:"first_name"`
	LastName        string          `gorm:"varchar(255);" faker:"last_name"`
	Mobile          uint64          `gorm:"not null;uniqueIndex"`
	MobileConfirmed bool            `gorm:"default:false"`
	Permissions     UserPermissions `gorm:"type:json;not null"`
	Password        string          `gorm:"varchar(255);not null"`
	Businesses      []*Business     `gorm:"many2many:business_users;" faker:"-"`
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
