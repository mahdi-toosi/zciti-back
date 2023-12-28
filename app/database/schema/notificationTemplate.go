package schema

import (
	"github.com/lib/pq"
)

type NotificationTemplate struct {
	ID      uint64         `gorm:"primaryKey" faker:"-"`
	Title   string         `gorm:"not null;varchar(250);uniqueIndex" faker:"sentence,unique"`
	Content string         `gorm:"not null" faker:"paragraph"`
	Tag     pq.StringArray `gorm:"type:text[]" faker:"oneof: tag,category,taxonomy"`
	Base
}