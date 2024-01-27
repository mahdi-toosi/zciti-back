package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Comment struct {
	ID       uint64
	Content  string `example:"content content content" validate:"required,max=1000"`
	Status   string
	AuthorID uint64
}

type UpdateCommentStatus struct {
	Status string `example:"pending" validate:"required,oneof=approved pending"`
}

type Comments struct {
	Pagination *paginator.Pagination
}

func (req *Comment) ToDomain(postID *uint64) *schema.Comment {
	return &schema.Comment{
		ID:       req.ID,
		Status:   req.Status,
		Content:  req.Content,
		AuthorID: req.AuthorID,
		PostID:   *postID,
	}
}
