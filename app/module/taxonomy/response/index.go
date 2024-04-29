package response

import (
	"go-fiber-starter/app/database/schema"
)

type Taxonomy struct {
	ID          uint64              `json:",omitempty"`
	Title       string              `json:",omitempty"`
	Description string              `json:",omitempty"`
	Type        schema.TaxonomyType `json:",omitempty"`
	ParentID    uint64              `json:",omitempty"`
}

func FromDomain(item *schema.Taxonomy) (res *Taxonomy) {
	if item != nil {
		res = &Taxonomy{
			ID:          item.ID,
			Type:        item.Type,
			Title:       item.Title,
			ParentID:    item.ParentID,
			Description: item.Description,
		}
	}

	return res
}
