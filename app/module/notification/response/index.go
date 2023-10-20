package response

import (
	"time"

	"go-fiber-starter/app/database/schema"
)

type Notification struct {
	ID              uint64
	FirstName       string
	LastName        string
	Mobile          uint64
	MobileConfirmed bool
	Roles           []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromDomain(notification *schema.Notification) (res *Notification) {
	if notification != nil {
		res = &Notification{
			ID: notification.ID,
			//FirstName:       notification.FirstName,
			//LastName:        notification.LastName,
			//Mobile:          notification.Mobile,
			//MobileConfirmed: notification.MobileConfirmed,
			//Roles:           notification.Roles,

			CreatedAt: notification.CreatedAt,
			UpdatedAt: notification.UpdatedAt,
		}
	}

	return res
}
