package response

import (
	ptime "github.com/yaa110/go-persian-calendar"
	"go-fiber-starter/app/database/schema"
	bresponse "go-fiber-starter/app/module/business/response"
	uresponse "go-fiber-starter/app/module/user/response"
	"time"
)

type Reservation struct {
	ID               uint64
	Status           schema.ReservationStatus
	EndTime          time.Time
	StartTime        time.Time
	EndTimeDisplay   string
	StartTimeDisplay string
	User             uresponse.User
	ProductID        uint64
	PostID           uint64
	PostTitle        string
	ProductSKU       string
	UserUsageCount   uint64
	Business         *bresponse.Business
	Meta             schema.ReservationMeta
}

func FromDomain(item *schema.Reservation) (res *Reservation) {
	if item == nil {
		return nil
	}

	res = &Reservation{
		ID:             item.ID,
		Meta:           item.Meta,
		Status:         item.Status,
		EndTime:        item.EndTime,
		StartTime:      item.StartTime,
		ProductID:      item.ProductID,
		UserUsageCount: item.UserUsageCount,
		PostID:         item.Product.Post.ID,
		ProductSKU:     item.Product.Meta.SKU,
		PostTitle:      item.Product.Post.Title,
	}

	if res.Meta.UniWashLastCommandTime == nil || res.Meta.UniWashLastCommandTime.IsZero() {
		res.Meta.UniWashLastCommandTime = nil
	}

	if item.User.ID != 0 {
		res.User = uresponse.User{
			ID:       item.User.ID,
			Mobile:   item.User.Mobile,
			FullName: item.User.FullName(),
		}
	}

	if item.Business.ID != 0 {
		res.Business = &bresponse.Business{
			ID:    item.Business.ID,
			Title: item.Business.Title,
		}
	}

	res.EndTimeDisplay = ptime.New(res.EndTime).Format("HH:mm - yyyy/MM/dd")
	res.StartTimeDisplay = ptime.New(res.StartTime).Format("HH:mm - yyyy/MM/dd")

	return res
}
