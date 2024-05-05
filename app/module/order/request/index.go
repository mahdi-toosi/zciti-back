package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"time"
)

type StoreUniWash struct {
	Date       string
	UserID     uint64
	PostID     uint64
	EndTime    string
	StartTime  string
	ProductID  uint64
	BusinessID uint64
}

func (s StoreUniWash) GetStartDateTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Tehran")
	startTime, _ := time.ParseInLocation(time.DateTime, s.Date+" "+s.StartTime, loc)
	return startTime.UTC()
}

func (s StoreUniWash) GetEndDateTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Tehran")
	endTime, _ := time.ParseInLocation(time.DateTime, s.Date+" "+s.EndTime, loc)
	return endTime.UTC()
}

type Order struct {
	ID            uint64
	Status        schema.OrderStatus        `example:"pending" validate:"omitempty,oneOf=pending processing onHold completed cancelled refunded failed"`
	PaymentMethod schema.OrderPaymentMethod `example:"page" validate:"required,oneof=product post page"`
	UserNote      string                    `example:"note note" validate:"omitempty,min=2,max=255" json:",omitempty" faker:""`
	BusinessID    uint64                    `example:"1" validate:"min=1"`
	UserID        uint64                    `example:"1" validate:"min=1"`
}

type Orders struct {
	BusinessID uint64
	Pagination *paginator.Pagination
}

func (req *Order) ToDomain() *schema.Order {
	return &schema.Order{
		ID:            req.ID,
		Status:        req.Status,
		UserID:        req.UserID,
		BusinessID:    req.BusinessID,
		PaymentMethod: req.PaymentMethod,
		Meta: schema.OrderMeta{
			UserNote: req.UserNote,
		},
	}
}
