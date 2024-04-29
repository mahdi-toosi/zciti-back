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

//func FromDomain(item *schema.Notification) (res *Notification) {
//	if item != nil {
//		res = &Notification{
//			ID:         item.ID,
//			Type:       item.Type,
//			BusinessID: item.BusinessID,
//			SentAt:     item.SentAt,
//			TemplateID: item.TemplateID,
//			//ReceiverID: item.ReceiverID,
//			Receiver: response.User{
//				ID:       item.Receiver.ID,
//				FullName: fmt.Sprint(item.Receiver.FirstName, " ", item.Receiver.LastName),
//			},
//
//			CreatedAt: item.CreatedAt,
//			UpdatedAt: item.UpdatedAt,
//		}
//	}
//
//	return res
//}
