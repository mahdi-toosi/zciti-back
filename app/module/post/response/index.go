package response

import (
	"go-fiber-starter/app/database/schema"
)

type Post struct {
	ID       uint64
	Title    string
	AuthorID uint64
	Content  string
	Status   string
	Type     string
}

func FromDomain(post *schema.Post) (res *Post) {
	if post != nil {
		res = &Post{
			ID:       post.ID,
			Title:    post.Title,
			AuthorID: post.AuthorID,
			Content:  post.Content,
			Type:     post.Type,
		}
	}

	return res
}
