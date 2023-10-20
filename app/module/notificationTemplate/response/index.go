package response

import (
	"encoding/json"
	"go-fiber-starter/app/database/schema"
)

type NotificationTemplate struct {
	ID      uint64
	Title   string
	Content string
	Meta    string
	Tag     []string
}

func FromDomain(nt *schema.NotificationTemplate) (res *NotificationTemplate) {
	meta := map[string]any{"mahdi": "toosi"}

	jsonMeta, _ := json.Marshal(meta)

	if nt != nil {
		res = &NotificationTemplate{
			ID:      nt.ID,
			Title:   nt.Title,
			Content: nt.Content,
			Meta:    string(jsonMeta),
			Tag:     nt.Tag,
		}
	}

	return res
}
