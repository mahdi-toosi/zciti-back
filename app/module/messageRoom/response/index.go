package response

import (
	"go-fiber-starter/app/database/schema"
)

type MessageRoom struct {
	ID         uint64 `json:",omitempty"`
	BusinessID uint64 `json:",omitempty"`
	UserID     uint64 `json:",omitempty"`
	Status     string `json:",omitempty"`
}

func FromDomain(messageRoom *schema.MessageRoom) (res *MessageRoom) {
	if messageRoom == nil {
		return nil
	}

	return &MessageRoom{
		ID:         messageRoom.ID,
		UserID:     messageRoom.UserID,
		Status:     messageRoom.Status,
		BusinessID: messageRoom.BusinessID,
	}
}
