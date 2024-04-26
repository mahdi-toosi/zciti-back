package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type User struct {
	ID          uint64
	FirstName   string `example:"mahdi" validate:"required,min=2,max=255"`
	LastName    string `example:"lastname" validate:"required,min=2,max=255"`
	Mobile      uint64 `example:"9380338494" validate:"required,number"`
	Password    string
	Permissions schema.UserPermissions `example:"{1:['operator']}"`
}

type BusinessUsers struct {
	Pagination *paginator.Pagination
	BusinessID uint64
	Keyword    string
}

type Users struct {
	Pagination *paginator.Pagination
	Keyword    string
}

func (req *User) ToDomain() *schema.User {
	return &schema.User{
		ID:          req.ID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Mobile:      req.Mobile,
		Permissions: req.Permissions,
	}
}
