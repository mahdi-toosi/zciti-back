package response

import (
	"go-fiber-starter/app/database/schema"
)

type NotificationTemplate struct {
	ID         uint64
	Title      string
	Content    string
	Tag        []string
	BusinessID uint64
}

func FromDomain(nt *schema.NotificationTemplate) (res *NotificationTemplate) {
	if nt != nil {
		res = &NotificationTemplate{
			ID:         nt.ID,
			Tag:        nt.Tag,
			Title:      nt.Title,
			Content:    nt.Content,
			BusinessID: nt.BusinessID,
		}
	}

	return res
}
