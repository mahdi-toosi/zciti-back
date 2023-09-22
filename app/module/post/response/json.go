package response

import (
	"time"

	"go-fiber-starter/app/database/schema"
)

type Post struct {
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	AuthorID uint64 `json:"author_id"`
	Content  string `json:"content"`
	Type     string `json:"type"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FromDomain(post *schema.Post) (res *Post) {
	if post != nil {
		res = &Post{
			ID:       post.ID,
			Title:    post.Title,
			AuthorID: post.AuthorID,
			Content:  post.Content,
			Type:     post.Type,

			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
	}

	return res
}
