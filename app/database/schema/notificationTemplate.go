package schema

import (
	"github.com/lib/pq"
)

type NotificationTemplate struct {
	ID         uint64         `gorm:"primaryKey" faker:"-"`
	Title      string         `gorm:"not null;varchar(255);" faker:"sentence,unique"`
	Content    string         `gorm:"not null" faker:"paragraph"`
	BusinessID uint64         `faker:"-"`
	Tag        pq.StringArray `gorm:"type:text[]" faker:"oneof: tag,category,taxonomy"`
	Base
}
