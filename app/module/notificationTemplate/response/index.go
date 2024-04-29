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

func FromDomain(item *schema.NotificationTemplate) (res *NotificationTemplate) {
	if item != nil {
		res = &NotificationTemplate{
			ID:         item.ID,
			Tag:        item.Tag,
			Title:      item.Title,
			Content:    item.Content,
			BusinessID: item.BusinessID,
		}
	}

	return res
}
