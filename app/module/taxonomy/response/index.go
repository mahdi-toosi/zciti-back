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

func FromDomain(taxonomy *schema.Taxonomy) (res *Taxonomy) {
	if taxonomy != nil {
		res = &Taxonomy{
			ID:          taxonomy.ID,
			Type:        taxonomy.Type,
			Title:       taxonomy.Title,
			ParentID:    taxonomy.ParentID,
			Description: taxonomy.Description,
		}
	}

	return res
}
