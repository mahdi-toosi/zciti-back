package response

import (
	"fmt"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/response"
)

type Business struct {
	ID          uint64
	Title       string
	Type        schema.BusinessType
	TypeDisplay string
	Description string
	Meta        string
	OwnerID     uint64
	Owner       response.User
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
		ID:    business.ID,
		Type:  business.Type,
		Title: business.Title,
		Owner: response.User{
			ID:       business.Owner.ID,
			FullName: fmt.Sprint(business.Owner.FirstName, " ", business.Owner.LastName),
		},
		Meta:        business.Meta,
		OwnerID:     business.OwnerID,
		Description: business.Description,
		TypeDisplay: schema.TypeDisplayProxy[business.Type],
	}
}
