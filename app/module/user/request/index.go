package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"strings"
)

type User struct {
	ID          uint64
	FirstName   string `example:"mahdi" validate:"required,min=2,max=255"`
	LastName    string `example:"lastname" validate:"required,min=2,max=255"`
	Mobile      uint64 `example:"9380338494" validate:"required,number"`
	Password    string
	Permissions schema.UserPermissions `example:"{1:['operator']}"`
}

type UpdateUserAccount struct {
	ID        uint64 `example:"1" validate:"required,min=1"`
	FirstName string `example:"mahdi" validate:"required,min=2,max=255"`
	LastName  string `example:"lastname" validate:"required,min=2,max=255"`
	Mobile    uint64 `example:"9380338494" validate:"required,number"`
}

type BusinessUsers struct {
	BusinessID uint64
	Username   string
	UserIDs    []uint64
	Pagination *paginator.Pagination
}

type BusinessUsersStoreRole struct {
	Roles      []schema.UserRole `example:"[user]" validate:"required"`
	UserID     uint64            `example:"1" validate:"required,number,min=1"`
	BusinessID uint64            `example:"1" validate:"required,number,min=1"`
}

type Users struct {
	Pagination *paginator.Pagination
	Keyword    string
}

func (req *User) ToDomain() *schema.User {
	return &schema.User{
		ID:          req.ID,
		Mobile:      req.Mobile,
		Permissions: req.Permissions,
		LastName:    strings.TrimSpace(req.LastName),
		FirstName:   strings.TrimSpace(req.FirstName),
	}
}

func (req *UpdateUserAccount) ToDomain() *schema.User {
	return &schema.User{
		ID:        req.ID,
		Mobile:    req.Mobile,
		LastName:  strings.TrimSpace(req.LastName),
		FirstName: strings.TrimSpace(req.FirstName),
	}
}
