package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type MessageRoom struct {
	ID         uint64
	BusinessID uint64 `example:"1" validate:"number"`
	UserID     uint64 `example:"1" validate:"number"`
	Status     string `example:"active" validate:"oneof:active,archived,blocked"`
}

type MessageRooms struct {
	Pagination *paginator.Pagination
}

func (req *MessageRoom) ToDomain() *schema.MessageRoom {
	return &schema.MessageRoom{
		ID:         req.ID,
		UserID:     req.UserID,
		Status:     req.Status,
		BusinessID: req.BusinessID,
	}
}
