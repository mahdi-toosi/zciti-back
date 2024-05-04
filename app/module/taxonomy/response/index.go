package response

import (
	"go-fiber-starter/app/database/schema"
)

type Taxonomy struct {
	ID          uint64              `json:",omitempty"`
	Title       string              `json:",omitempty"`
	Description string              `json:",omitempty"`
	ParentID    *uint64             `json:",omitempty"`
	Domain      schema.PostType     `json:",omitempty"`
	Type        schema.TaxonomyType `json:",omitempty"`
}

func FromDomain(item *schema.Taxonomy, forUser bool) (res *Taxonomy) {
	if item == nil {
		return res
	}

	if forUser {
		return &Taxonomy{
			ID:       item.ID,
			Type:     item.Type,
			Title:    item.Title,
			Domain:   item.Domain,
			ParentID: item.ParentID,
		}
	}

	res = &Taxonomy{
		ID:          item.ID,
		Type:        item.Type,
		Title:       item.Title,
		Domain:      item.Domain,
		ParentID:    item.ParentID,
		Description: item.Description,
	}

	return res
}
