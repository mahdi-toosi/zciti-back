package response

import (
	"time"
)

type Comment struct {
	ID              uint64 `json:",omitempty"`
	Content         string `json:",omitempty"`
	Status          string `json:",omitempty"`
	AuthorID        uint64 `json:",omitempty"`
	PostID          uint64 `json:",omitempty"`
	ParentID        uint64 `json:",omitempty"`
	IsBusinessOwner bool
	AuthorFullName  string    `json:",omitempty"`
	CreatedAt       time.Time `json:",omitempty"`
}

func FromDomain(comment map[string]interface{}) (res *Comment) {
	if comment != nil {
		res = &Comment{
			ID:              uint64(comment["id"].(int64)),
			Content:         comment["content"].(string),
			Status:          comment["status"].(string),
			CreatedAt:       comment["created_at"].(time.Time),
			IsBusinessOwner: comment["is_business_owner"].(bool),
			AuthorID:        uint64(comment["author_id"].(int64)),
			AuthorFullName:  comment["author_full_name"].(string),
		}

		if rec, ok := comment["parent_id"].(int64); ok {
			res.ParentID = uint64(rec)
		}
	}

	return res
}
