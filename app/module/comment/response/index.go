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

func FromDomain(item map[string]interface{}) (res *Comment) {
	if item != nil {
		res = &Comment{
			ID:              uint64(item["id"].(int64)),
			Content:         item["content"].(string),
			Status:          item["status"].(string),
			CreatedAt:       item["created_at"].(time.Time),
			IsBusinessOwner: item["is_business_owner"].(bool),
			AuthorID:        uint64(item["author_id"].(int64)),
			AuthorFullName:  item["author_full_name"].(string),
		}

		if rec, ok := item["parent_id"].(int64); ok {
			res.ParentID = uint64(rec)
		}
	}

	return res
}
