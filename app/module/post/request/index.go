package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Post struct {
	ID             uint64
	AuthorID       uint64 `example:"1" validate:"number"`
	BusinessID     uint64 `example:"1" validate:"number"`
	Title          string `example:"title" validate:"min=2,max=255"`
	Content        string `example:"content content content" validate:"min=2,max=255"`
	Status         string `example:"draft" validate:"oneof=draft published private"`
	Type           string `example:"page" validate:"oneof=product post page"`
	CommentsStatus string `example:"open" validate:"oneof=open close onlyBuyers onlyCustomers"`
}

type PostsRequest struct {
	BusinessID uint64
	Pagination *paginator.Pagination
}

func (req *Post) ToDomain() *schema.Post {
	return &schema.Post{
		ID:             req.ID,
		Type:           req.Type,
		Title:          req.Title,
		Status:         req.Status,
		Content:        req.Content,
		AuthorID:       req.AuthorID,
		BusinessID:     req.BusinessID,
		CommentsStatus: req.CommentsStatus,
	}
}
