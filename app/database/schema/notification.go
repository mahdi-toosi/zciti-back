package schema

import (
	"github.com/lib/pq"
	"time"
)

type Notification struct {
	ID         uint64         `gorm:"primaryKey" faker:"-"`
	ReceiverID uint64         `gorm:"not null" faker:"-"`
	Receiver   User           `gorm:"foreignKey:ReceiverID" faker:"-"`
	Type       pq.StringArray `gorm:"type:text[];not null" faker:"-"`
	BusinessID uint64         `gorm:"not null" faker:"-"`
	SentAt     time.Time      `gorm:"not null"`
	TemplateID uint64         `faker:"-"`
	//Template   NotificationTemplate `gorm:"foreignKey:TemplateID" faker:"-"`
	Base
}

type NotificationType string

const (
	TSms          NotificationType = "Sms"
	TNotification NotificationType = "Notification"
)
