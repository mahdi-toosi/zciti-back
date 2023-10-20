package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Post struct {
	ID       uint64 `json:"id"`
	AuthorID uint64 `json:"author_id" example:"1"`
	Title    string `json:"title" example:"title" validate:"min=2,max=255"`
	Content  string `json:"content" example:"content content content" validate:"min=2,max=255"`
	Status   string `json:"status" example:"draft"`
	Type     string `json:"type" example:"page"`
}

type PostsRequest struct {
	Pagination *paginator.Pagination
}

func (req *Post) ToDomain() *schema.Post {
	return &schema.Post{
		ID:       req.ID,
		Title:    req.Title,
		AuthorID: req.AuthorID,
		Content:  req.Content,
		Type:     req.Type,
	}
}
