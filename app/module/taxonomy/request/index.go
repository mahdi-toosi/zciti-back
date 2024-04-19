package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Taxonomy struct {
	ID          uint64
	Title       string              `example:"title" validate:"required,min=2,max=100"`
	Description string              `example:"description" validate:"min=2,max=500"`
	Type        schema.TaxonomyType `example:"tag" validate:"required,oneof=category tag"`
	Domain      schema.PostType     `example:"post" validate:"required,oneof=post page product"`
	BusinessID  uint64              `example:"1" validate:"number"`
	Business    schema.Business     ``
	ParentID    uint64              `example:"1" validate:"number"`
}

type Taxonomies struct {
	BusinessID uint64
	Keyword    string
	Domain     schema.PostType     `example:"post" validate:"omitempty,oneof=post page product"`
	Type       schema.TaxonomyType `example:"post" validate:"omitempty,oneof=category tag"`
	Pagination *paginator.Pagination
}

func (req *Taxonomy) ToDomain() *schema.Taxonomy {
	return &schema.Taxonomy{
		ID:          req.ID,
		Type:        req.Type,
		Title:       req.Title,
		Business:    req.Business,
		ParentID:    req.ParentID,
		Description: req.Description,
	}
}
