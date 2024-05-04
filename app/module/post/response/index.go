package response

import (
	"go-fiber-starter/app/database/schema"
	tresponse "go-fiber-starter/app/module/taxonomy/response"
	"go-fiber-starter/app/module/user/response"
)

type Post struct {
	ID         uint64            `json:",omitempty"`
	Title      string            `json:",omitempty"`
	AuthorID   uint64            `json:",omitempty"`
	Excerpt    string            `json:",omitempty"`
	Author     response.User     `json:",omitempty"`
	Content    string            `json:",omitempty"`
	Status     schema.PostStatus `json:",omitempty"`
	Type       schema.PostType   `json:",omitempty"`
	BusinessID uint64            `json:",omitempty"`
	//Business   bresponse.Business   `json:",omitempty"`
	Taxonomies []tresponse.Taxonomy `json:",omitempty"`
	Meta       schema.PostMeta      `json:",omitempty"`
}

func FromDomain(item *schema.Post) (res *Post) {
	if item == nil {
		return res
	}

	p := &Post{
		ID:      item.ID,
		Type:    item.Type,
		Meta:    item.Meta,
		Title:   item.Title,
		Status:  item.Status,
		Content: item.Content,
		Author:  response.User{ID: item.Author.ID, FullName: item.Author.FullName()},
		//Business: bresponse.Business{ID: item.Business.ID, Title: item.Business.Title},
	}

	for _, taxonomy := range item.Taxonomies {
		p.Taxonomies = append(p.Taxonomies, tresponse.Taxonomy{
			ID:       taxonomy.ID,
			Type:     taxonomy.Type,
			Title:    taxonomy.Title,
			ParentID: taxonomy.ParentID,
		})
	}

	return p
}
