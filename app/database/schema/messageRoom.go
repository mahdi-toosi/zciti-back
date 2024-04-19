package schema

type MessageRoom struct {
	ID         uint64 `gorm:"primaryKey" faker:"-"`
	BusinessID uint64 `gorm:"not null; index:idx_room;priority:1" faker:"-"`
	UserID     uint64 `gorm:"not null; index:idx_room;priority:2" faker:"-"`
	Status     string `gorm:"varchar(250);" faker:"oneof:active,archived,blocked"`
	Base
}

type MessageRoomStatus string

const (
	MessageRoomStatusActive   MessageRoomStatus = "active"
	MessageRoomStatusBlocked  MessageRoomStatus = "blocked"
	MessageRoomStatusArchived MessageRoomStatus = "archived"
)
