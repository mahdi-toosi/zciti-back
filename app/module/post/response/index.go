package response

import (
	ptime "github.com/yaa110/go-persian-calendar"
	"go-fiber-starter/app/database/schema"
	tresponse "go-fiber-starter/app/module/taxonomy/response"
	"go-fiber-starter/app/module/user/response"
	"time"
)

type Post struct {
	ID               uint64               `json:",omitempty"`
	Title            string               `json:",omitempty"`
	AuthorID         uint64               `json:",omitempty"`
	Excerpt          string               `json:",omitempty"`
	Author           response.User        `json:",omitempty"`
	Content          string               `json:",omitempty"`
	Status           schema.PostStatus    `json:",omitempty"`
	Type             schema.PostType      `json:",omitempty"`
	BusinessID       uint64               `json:",omitempty"`
	Taxonomies       []tresponse.Taxonomy `json:",omitempty"`
	Meta             schema.PostMeta      `json:",omitempty"`
	CreatedAt        time.Time            `json:",omitempty"`
	CreatedAtDisplay string               `json:",omitempty"`
	Observers        []*response.User     ``
	//Business   bresponse.Business   `json:",omitempty"`
}

func FromDomain(item *schema.Post) (res *Post) {
	if item == nil {
		return res
	}

	p := &Post{
		ID:               item.ID,
		Type:             item.Type,
		Meta:             item.Meta,
		Title:            item.Title,
		Status:           item.Status,
		Content:          item.Content,
		Excerpt:          item.Excerpt,
		CreatedAt:        item.CreatedAt,
		CreatedAtDisplay: ptime.New(item.CreatedAt).Format("HH:mm - MM/dd"),
		Author:           response.User{ID: item.Author.ID, FullName: item.Author.FullName()},
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
