package response

import (
	"fmt"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/response"
	"time"
)

type Comment struct {
	ID        uint64        `json:",omitempty"`
	Content   string        `json:",omitempty"`
	Status    string        `json:",omitempty"`
	AuthorID  uint64        `json:",omitempty"`
	PostID    uint64        `json:",omitempty"`
	Author    response.User `json:",omitempty"`
	CreatedAt time.Time     `json:",omitempty"`
}

func FromDomain(comment *schema.Comment) (res *Comment) {
	if comment != nil {
		res = &Comment{
			ID:        comment.ID,
			Content:   comment.Content,
			Status:    comment.Status,
			CreatedAt: comment.CreatedAt,
			Author: response.User{
				ID:       comment.Author.ID,
				FullName: fmt.Sprint(comment.Author.FirstName, " ", comment.Author.LastName),
			},
		}
	}

	return res
}
