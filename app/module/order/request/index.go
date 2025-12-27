package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/orderItem/request"
	"go-fiber-starter/utils/paginator"
	"time"
)

type Order struct {
	ID            uint64
	Status        schema.OrderStatus        `example:"pending" validate:"omitempty,oneof=pending processing onHold completed cancelled refunded failed"`
	PaymentMethod schema.OrderPaymentMethod `example:"online" validate:"required,oneof=cash online cashOnDelivery"`
	UserNote      string                    `example:"note note" validate:"omitempty,min=2,max=255" json:",omitempty" faker:""`
	BusinessID    uint64                    `example:"1" validate:"min=1"`
	CouponCode    string                    `example:"code"`
	CouponID      *uint64
	User          schema.User
	OrderItems    []request.OrderItem
}

type Orders struct {
	BusinessID     uint64
	CouponID       uint64
	UserID         uint64
	ProductID      uint64
	CityID         uint64
	WorkspaceID    uint64
	DormitoryID    uint64
	Taxonomies     []uint64
	Status         string     `example:"pending" validate:"omitempty,oneof=pending processing onHold completed cancelled refunded failed"`
	FullName       string     // Filter by user full name
	HasCoupon      *bool      // Filter by whether order has coupon (true = has coupon, false = no coupon)
	StartTime      *time.Time // Filter by reservation start time
	EndTime        *time.Time // Filter by reservation end time
	OrderStartTime *time.Time // Filter by order creation start time
	OrderEndTime   *time.Time // Filter by order creation end time
	Pagination     *paginator.Pagination
}

func (req *Order) ToDomain(totalAmt *float64, authority *string) *schema.Order {
	o := &schema.Order{
		ID:            req.ID,
		Status:        req.Status,
		UserID:        req.User.ID,
		BusinessID:    req.BusinessID,
		PaymentMethod: req.PaymentMethod,
		Meta: schema.OrderMeta{
			UserNote: req.UserNote,
		},
	}

	if totalAmt != nil {
		o.TotalAmt = *totalAmt
		o.Meta.TaxAmt = uint64(*totalAmt * 0.1)
	}

	if req.CouponID != nil {
		o.CouponID = req.CouponID
	}

	if totalAmt != nil && int(*totalAmt) == 0 {
		req.Status = schema.OrderStatusCompleted
	}

	if authority != nil {
		o.Meta.PaymentAuthority = *authority
	}

	return o
}
