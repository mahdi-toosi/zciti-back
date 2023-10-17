package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type UserRequest struct {
	ID        uint64
	FirstName string   `example:"mahdi" validate:"min=2,max=255"`
	LastName  string   `example:"lastname" validate:"min=2,max=255"`
	Mobile    uint64   `example:"9150338494" validate:"number"`
	Roles     []string `example:"user"`
}

type UsersRequest struct {
	Pagination *paginator.Pagination `json:"pagination"`
}

func (req *UserRequest) ToDomain() *schema.User {
	return &schema.User{
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
		Roles:     req.Roles,
	}
}
