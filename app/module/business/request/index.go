package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Business struct {
	ID          uint64
	Title       string              `example:"title" validate:"required,min=2,max=255"`
	Type        schema.BusinessType `example:"type" validate:"required,min=2,max=255"`
	Description string              `example:"business" validate:"max=500"`
	OwnerID     uint64              `example:"1" validate:"required,number"`
}

type Businesses struct {
	Pagination *paginator.Pagination
	Keyword    string
}

type Users struct {
	Pagination *paginator.Pagination
	BusinessID uint64
	Keyword    string
}

func (req *Business) ToDomain() *schema.Business {
	return &schema.Business{
		ID:          req.ID,
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		OwnerID:     req.OwnerID,
	}
}
