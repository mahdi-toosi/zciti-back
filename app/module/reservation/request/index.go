package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"time"
)

type Reservation struct {
	ID         uint64
	ReceiverID uint64    `reservation:"1" validate:"required,number,min=1"`
	Type       []string  `reservation:"Sms" validate:"required,dive"`
	BusinessID uint64    `reservation:"1" validate:"min=1"`
	SentAt     time.Time `reservation:"2023-10-20T15:47:33.084Z" validate:"datetime"`
	TemplateID uint64    `reservation:"1" validate:"required,min=1"`
}

type Reservations struct {
	BusinessID uint64
	UserID     uint64
	Mobile     string
	FullName   string
	ProductID  uint64
	StartTime  *time.Time
	EndTime    *time.Time
	Status     schema.ReservationStatus
	Pagination *paginator.Pagination
}

func (req *Reservation) ToDomain() *schema.Reservation {
	return &schema.Reservation{
		//ID:         req.ID,
		//ReceiverID: req.ReceiverID,
		//Type:       req.Type,
		//BusinessID: req.BusinessID,
		//SentAt:     req.SentAt,
		//TemplateID: req.TemplateID,
	}
}
