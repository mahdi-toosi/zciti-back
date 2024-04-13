package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Taxonomy struct {
	ID          uint64
	FirstName   string                 `example:"mahdi" validate:"required,min=2,max=255"`
	LastName    string                 `example:"lastname" validate:"required,min=2,max=255"`
	Mobile      uint64                 `example:"9380338494" validate:"required,number"`
	Permissions schema.UserPermissions `example:"taxonomy"`
}

type Taxonomies struct {
	Pagination *paginator.Pagination
	Keyword    string
}

func (req *Taxonomy) ToDomain() *schema.Taxonomy {
	return &schema.Taxonomy{
		ID:          req.ID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Mobile:      req.Mobile,
		Permissions: req.Permissions,
	}
}
