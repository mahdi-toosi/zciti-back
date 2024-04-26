package response

import (
	"go-fiber-starter/app/database/schema"
	bresponse "go-fiber-starter/app/module/business/response"
	tresponse "go-fiber-starter/app/module/taxonomy/response"
	"go-fiber-starter/app/module/user/response"
)

type Post struct {
	ID         uint64               `json:",omitempty"`
	Title      string               `json:",omitempty"`
	AuthorID   uint64               `json:",omitempty"`
	Author     response.User        `json:",omitempty"`
	Content    string               `json:",omitempty"`
	Status     schema.PostStatus    `json:",omitempty"`
	Type       schema.PostType      `json:",omitempty"`
	BusinessID uint64               `json:",omitempty"`
	Business   bresponse.Business   `json:",omitempty"`
	Taxonomies []tresponse.Taxonomy `json:",omitempty"`
	Meta       schema.PostMeta      `json:",omitempty"`
}

func FromDomain(post *schema.Post) (res *Post) {
	if post == nil {
		return res
	}

	p := &Post{
		ID:       post.ID,
		Type:     post.Type,
		Meta:     post.Meta,
		Title:    post.Title,
		Status:   post.Status,
		Content:  post.Content,
		Author:   response.User{ID: post.Author.ID, FullName: post.Author.FullName()},
		Business: bresponse.Business{ID: post.Business.ID, Title: post.Business.Title},
	}

	for _, taxonomy := range post.Taxonomies {
		p.Taxonomies = append(p.Taxonomies, tresponse.Taxonomy{
			ID:       taxonomy.ID,
			Type:     taxonomy.Type,
			Title:    taxonomy.Title,
			ParentID: taxonomy.ParentID,
		})
	}

	return p
}
