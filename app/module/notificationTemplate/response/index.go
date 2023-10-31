package response

import (
	"go-fiber-starter/app/database/schema"
)

type NotificationTemplate struct {
	ID      uint64
	Title   string
	Content string
	Tag     []string
}

func FromDomain(nt *schema.NotificationTemplate) (res *NotificationTemplate) {
	if nt != nil {
		res = &NotificationTemplate{
			ID:      nt.ID,
			Title:   nt.Title,
			Content: nt.Content,
			Tag:     nt.Tag,
		}
	}

	return res
}
