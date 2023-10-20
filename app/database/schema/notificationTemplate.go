package schema

import "github.com/lib/pq"

type NotificationTemplate struct {
	ID      uint64         `gorm:"primary_key" faker:"-"`
	Title   string         `gorm:"not null;varchar(250);unique" faker:"sentence"`
	Content string         `gorm:"not null" faker:"paragraph"`
	Meta    string         `gorm:"default:'[]'" faker:"-"`
	Tag     pq.StringArray `gorm:"type:text[]" faker:"oneof: tag,category,taxonomy"`
	Base
}
