package response

import (
	"go-fiber-starter/app/database/schema"
	"time"
)

type Message struct {
	ID        uint64    `json:",omitempty"`
	RoomID    uint64    `json:",omitempty"`
	FromID    uint64    `json:",omitempty"`
	ToID      uint64    `json:",omitempty"`
	Content   string    `json:",omitempty"`
	CreatedAt time.Time `json:",omitempty"`
}

func FromDomain(message *schema.Message) (res *Message) {
	if message == nil {
		return nil
	}

	return &Message{
		ID:        message.ID,
		RoomID:    message.RoomID,
		FromID:    message.FromID,
		ToID:      message.ToID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}
}
