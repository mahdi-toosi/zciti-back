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

func FromDomain(item *schema.Message) (res *Message) {
	if item == nil {
		return nil
	}

	return &Message{
		ID:        item.ID,
		RoomID:    item.RoomID,
		FromID:    item.FromID,
		ToID:      item.ToID,
		Content:   item.Content,
		CreatedAt: item.CreatedAt,
	}
}
