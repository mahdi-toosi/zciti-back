package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type PostRequest struct {
	ID       uint64 `json:"id"`
	AuthorID uint64 `json:"author_id"`
	Title    string `json:"title" validate:"min=2,max=255"`
	Content  string `json:"content" validate:"min=2,max=255"`
	Status   string `json:"status"`
	Type     string `json:"type"`
}

type PostsRequest struct {
	Pagination *paginator.Pagination `json:"pagination"`
}

func (req *PostRequest) ToDomain() *schema.Post {
	return &schema.Post{
		ID:       req.ID,
		Title:    req.Title,
		AuthorID: req.AuthorID,
		Content:  req.Content,
		Type:     req.Type,
	}
}
