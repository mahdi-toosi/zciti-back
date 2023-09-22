package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type UserRequest struct {
	ID        uint64   `json:"id"`
	FirstName string   `json:"firstName" validate:"min=2,max=255"`
	LastName  string   `json:"lastName" validate:"min=2,max=255"`
	Mobile    uint64   `json:"mobile" validate:"regex:09(1[0-9]|3[1-9]|2[1-9])-?[0-9]{3}-?[0-9]{4}"`
	Roles     []string `json:"roles"`
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
