package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type NotificationTemplate struct {
	ID         uint64
	Title      string   `example:"title" validate:"min=2,max=255"`
	Content    string   `example:"some content some content some content" validate:"min=2"`
	Tag        []string `example:"['tag']" validate:"dive"`
	BusinessID uint64
}

type Index struct {
	BusinessID uint64
	Pagination *paginator.Pagination
}

func (req *NotificationTemplate) ToDomain() *schema.NotificationTemplate {
	return &schema.NotificationTemplate{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Tag:     req.Tag,
	}
}
