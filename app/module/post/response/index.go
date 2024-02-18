package response

import (
	"go-fiber-starter/app/database/schema"
	bresponse "go-fiber-starter/app/module/business/response"
	"go-fiber-starter/app/module/user/response"
)

type Post struct {
	ID             uint64             `json:",omitempty"`
	Title          string             `json:",omitempty"`
	AuthorID       uint64             `json:",omitempty"`
	Author         response.User      `json:",omitempty"`
	Content        string             `json:",omitempty"`
	Status         string             `json:",omitempty"`
	Type           string             `json:",omitempty"`
	BusinessID     uint64             `json:",omitempty"`
	Business       bresponse.Business `json:",omitempty"`
	CommentsStatus string             `json:",omitempty"`
	CommentsCount  uint64             `json:",omitempty"`
}

func FromDomain(post *schema.Post) (res *Post) {
	if post != nil {
		res = &Post{
			ID:             post.ID,
			Title:          post.Title,
			Author:         response.User{ID: post.Author.ID, FullName: post.Author.FirstName + " " + post.Author.LastName},
			Content:        post.Content,
			Type:           post.Type,
			Status:         post.Status,
			Business:       bresponse.Business{ID: post.Business.ID, Title: post.Business.Title},
			CommentsStatus: post.CommentsStatus,
			CommentsCount:  post.CommentsCount,
		}
	}

	return res
}
