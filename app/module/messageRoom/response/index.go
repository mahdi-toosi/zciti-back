package response

import (
	"go-fiber-starter/app/database/schema"
)

type MessageRoom struct {
	RoomID     uint64  `json:",omitempty"`
	BusinessID uint64  `json:",omitempty"`
	UserID     uint64  `json:",omitempty"`
	Status     string  `json:",omitempty"`
	Token      *string `json:",omitempty"`
}

type MessageRoomToken struct {
	//schema.MessageRoom
	Members         []string
	MembersAsString string
	ID              uint64
	BusinessID      uint64
	UserID          uint64
	Status          string
}

func FromDomain(item *schema.MessageRoom, token *string) (res *MessageRoom) {
	if item == nil {
		return nil
	}

	return &MessageRoom{
		Token:      token,
		RoomID:     item.ID,
		UserID:     item.UserID,
		Status:     item.Status,
		BusinessID: item.BusinessID,
	}
}
