package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"time"
)

type Notification struct {
	ID         uint64
	ReceiverID uint64    `example:"1" validate:"required,number,min=1"`
	Type       []string  `example:"Sms" validate:"required,dive"`
	BusinessID uint64    `example:"1" validate:"min=1"`
	SentAt     time.Time `example:"2023-10-20T15:47:33.084Z" validate:"datetime"`
	TemplateID uint64    `example:"1" validate:"required,min=1"`
}

type Notifications struct {
	BusinessID uint64
	Pagination *paginator.Pagination
}

func (req *Notification) ToDomain() *schema.Notification {
	return &schema.Notification{
		ID:         req.ID,
		ReceiverID: req.ReceiverID,
		Type:       req.Type,
		BusinessID: req.BusinessID,
		SentAt:     req.SentAt,
		TemplateID: req.TemplateID,
	}
}
