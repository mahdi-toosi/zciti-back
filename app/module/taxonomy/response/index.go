package response

import (
	"go-fiber-starter/app/database/schema"
)

type Taxonomy struct {
	ID              uint64                 `json:",omitempty"`
	FirstName       string                 `json:",omitempty"`
	LastName        string                 `json:",omitempty"`
	FullName        string                 `json:",omitempty"`
	Mobile          uint64                 `json:",omitempty"`
	MobileConfirmed bool                   `json:",omitempty"`
	Permissions     schema.UserPermissions `json:",omitempty"`
}

func FromDomain(taxonomy *schema.Taxonomy) (res *Taxonomy) {
	if taxonomy != nil {
		res = &Taxonomy{
			ID:              taxonomy.ID,
			FullName:        taxonomy.FullName(),
			Mobile:          taxonomy.Mobile,
			MobileConfirmed: taxonomy.MobileConfirmed,
			Permissions:     taxonomy.Permissions,
		}
	}

	return res
}
