package response

import (
	"time"

	"go-fiber-starter/app/database/schema"
)

type Post struct {
	ID              uint64   `json:"id"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	Mobile          string   `json:"mobile"`
	MobileConfirmed bool     `json:"mobileConfirmed"`
	Roles           []string `json:"roles"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FromDomain(post *schema.Post) (res *Post) {
	if post != nil {
		res = &Post{
			ID:              post.ID,
			FirstName:       post.FirstName,
			LastName:        post.LastName,
			Mobile:          post.Mobile,
			MobileConfirmed: post.MobileConfirmed,

			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
	}

	return res
}
