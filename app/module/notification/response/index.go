package response

import (
	"github.com/lib/pq"
	"time"
)

type Notification struct {
	ID               uint64
	ReceiverID       uint64
	ReceiverFullName string
	Type             pq.StringArray
	BusinessID       uint64
	BusinessTitle    string
	SentAt           time.Time
	TemplateID       uint64

	CreatedAt time.Time
	UpdatedAt time.Time
}

//func FromDomain(notification *schema.Notification) (res *Notification) {
//	if notification != nil {
//		res = &Notification{
//			ID:         notification.ID,
//			Type:       notification.Type,
//			BusinessID: notification.BusinessID,
//			SentAt:     notification.SentAt,
//			TemplateID: notification.TemplateID,
//			//ReceiverID: notification.ReceiverID,
//			Receiver: response.User{
//				ID:       notification.Receiver.ID,
//				FullName: fmt.Sprint(notification.Receiver.FirstName, " ", notification.Receiver.LastName),
//			},
//
//			CreatedAt: notification.CreatedAt,
//			UpdatedAt: notification.UpdatedAt,
//		}
//	}
//
//	return res
//}
