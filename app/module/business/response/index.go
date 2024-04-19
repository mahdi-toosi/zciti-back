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
	Owner       response.User       `json:",omitempty"`
	Meta        schema.BusinessMeta `json:",omitempty"`
}

type BusinessTypesOption struct {
	Label string
	Value schema.BusinessType
}

func FromDomain(business *schema.Business) (res *Business) {
	if business == nil {
		return res
	}

	return &Business{
		ID:          business.ID,
		Type:        business.Type,
		Meta:        business.Meta,
		Title:       business.Title,
		OwnerID:     business.OwnerID,
		Description: business.Description,
		TypeDisplay: schema.TypeDisplayProxy[business.Type],
		Owner: response.User{
			ID:       business.Owner.ID,
			FullName: business.Owner.FullName(),
		},
	}
}
