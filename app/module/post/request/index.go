package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Post struct {
	ID             uint64
	AuthorID       uint64                   `example:"1" validate:"number"`
	BusinessID     uint64                   `example:"1" validate:"number"`
	Title          string                   `example:"title" validate:"required,min=2,max=255"`
	Content        string                   `example:"content content content" validate:"required,min=2,max=255"`
	Status         schema.PostStatus        `example:"draft" validate:"required,oneof=draft published private"`
	Type           schema.PostType          `example:"page" validate:"required,oneof=product post page"`
	CommentsStatus schema.PostCommentStatus `example:"open" validate:"required,oneof=open close onlyBuyers onlyCustomers"`
}

type PostTaxonomies struct {
	BusinessID uint64
	PostID     uint64
	IDs        []uint64 `example:"[1,2,3]" validate:"dive"`
}

type PostsRequest struct {
	BusinessID uint64
	Keyword    string
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
