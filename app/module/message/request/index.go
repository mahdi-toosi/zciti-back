package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
)

type Message struct {
	ID      uint64
	RoomID  uint64 `example:"1" validate:"required,number"`
	FromID  uint64 `example:"1" validate:"required,number"`
	ToID    uint64 `example:"1" validate:"number"`
	Type    string `example:"text" validate:"required,oneof=text image"`
	Content string `example:"bla bla bla" validate:"required,max:2000"`
}

type Messages struct {
	BusinessID uint64
	UserID     uint64
	RoomID     *uint64
	Pagination *paginator.Pagination
}

func (req *Message) ToDomain() *schema.Message {
	return &schema.Message{
		ID:      req.ID,
		RoomID:  req.RoomID,
		FromID:  req.FromID,
		ToID:    req.ToID,
		Type:    req.Type,
		Content: req.Content,
	}
}
