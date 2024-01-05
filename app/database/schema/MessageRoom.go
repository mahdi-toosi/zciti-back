package schema

type MessageRoom struct {
	ID         uint64 `gorm:"primaryKey" faker:"-"`
	BusinessID uint64 `gorm:"not null; index:idx_room;" faker:"-"`
	UserID     uint64 `gorm:"not null; index:idx_room;" faker:"-"`
	Status     string `gorm:"varchar(250);" faker:"oneof:active,archived,blocked"`
	Base
}
