package schema

import "time"

type Message struct {
	ID      uint64    `gorm:"primaryKey" faker:"-"`
	RoomID  uint64    `gorm:"not null; index;" faker:"-"`
	FromID  uint64    `gorm:"not null" faker:"-"`
	ToID    uint64    `gorm:"not null" faker:"-"`
	Type    string    `gorm:"not null; varchar(100);" faker:"oneof:text,image,voice"`
	Content string    `gorm:"not null; varchar(2000);" faker:"paragraph"`
	SeenAt  time.Time `faker:"-"`
	Base
}
