package response

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/response"
)

type Business struct {
	ID          uint64              `json:",omitempty"`
	Title       string              `json:",omitempty"`
	Type        schema.BusinessType `json:",omitempty"`
	TypeDisplay string              `json:",omitempty"`
	Description string              `json:",omitempty"`
	OwnerID     uint64              `json:",omitempty"`
	Owner       *response.User      `json:",omitempty"`
	Meta        schema.BusinessMeta `json:",omitempty"`
}

type BusinessTypesOption struct {
	Label string
	Value schema.BusinessType
}

func FromDomain(item *schema.Business) (res *Business) {
	if item == nil {
		return res
	}

	b := &Business{
		ID:   item.ID,
		Type: item.Type,
		//Meta:        item.Meta,
		Title:       item.Title,
		OwnerID:     item.OwnerID,
		Description: item.Description,
		TypeDisplay: schema.TypeDisplayProxy[item.Type],
	}

	if item.Owner.ID != 0 {
		b.Owner = &response.User{
			ID:       item.Owner.ID,
			FullName: item.Owner.FullName(),
		}
	}

	return b
}
