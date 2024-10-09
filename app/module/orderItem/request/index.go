package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"time"
)

type OrderItems struct {
	Pagination *paginator.Pagination
}

type OrderItem struct {
	ID        uint64
	Quantity  int    `example:"1" validate:"number,min=1"`
	PostID    uint64 `example:"1" validate:"number,min=1"`
	ProductID uint64 `example:"1" validate:"number,min=1"`
	OrderID   uint64
	Date      string // for reservable
	EndTime   string // for reservable
	StartTime string // for reservable
}

func (req *OrderItem) GetStartDateTime() time.Time {
	startTime, _ := time.Parse(time.DateTime, req.Date+" "+req.StartTime)
	return startTime.UTC()
}

func (req *OrderItem) GetEndDateTime() time.Time {
	endTime, _ := time.Parse(time.DateTime, req.Date+" "+req.EndTime)
	return endTime.UTC()
}

type ToDomainParams struct {
	PostID        uint64
	Quantity      int
	Post          schema.Post
	Product       schema.Product
	ReservationID *uint64
}

func ToDomain(p ToDomainParams) *schema.OrderItem {
	orderItem := schema.OrderItem{
		PostID:        p.PostID,
		Quantity:      p.Quantity,
		ReservationID: p.ReservationID,
		Price:         p.Product.Price,
		Type:          schema.OrderItemTypeReservation,
		Subtotal:      float64(p.Quantity) * p.Product.Price,
		Meta: schema.OrderItemMeta{
			ProductTitle:       p.Post.Title,
			ProductID:          p.Product.ID,
			ProductType:        p.Product.Type,
			ProductSKU:         p.Product.Meta.SKU,
			ProductDetail:      p.Product.Meta.Detail,
			ProductVariantType: *p.Product.VariantType,
			//ProductImage:  post.Image,
		},
	}

	return &orderItem
}
