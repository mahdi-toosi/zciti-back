package response

import (
	ptime "github.com/yaa110/go-persian-calendar"
	"go-fiber-starter/app/database/schema"
	tresponse "go-fiber-starter/app/module/taxonomy/response"
	"time"
)

type Reservation struct {
	ID               uint64                `json:",omitempty"`
	EndTime          time.Time             `json:",omitempty"`
	StartTime        time.Time             `json:",omitempty"`
	EndTimeDisplay   string                `json:",omitempty"`
	StartTimeDisplay string                `json:",omitempty"`
	ProductID        uint64                `json:",omitempty"`
	ProductSKU       string                `json:",omitempty"`
	LastCommand      schema.UniWashCommand `json:",omitempty"`
	ProductTitle     string                `json:",omitempty"`
	ProductDetail    string                `json:",omitempty"`
}

func FromDomain(item *schema.Reservation) (res *Reservation) {
	if item == nil {
		return res
	}

	p := &Reservation{
		ID:            item.ID,
		EndTime:       item.EndTime,
		StartTime:     item.StartTime,
		ProductID:     item.ProductID,
		ProductSKU:    item.Product.Meta.SKU,
		ProductTitle:  item.Product.Post.Title,
		ProductDetail: item.Product.Meta.Detail,
		LastCommand:   item.Meta.UniWashLastCommand,
	}

	p.EndTimeDisplay = ptime.New(p.EndTime).Format("HH:mm - yyyy/MM/dd")
	p.StartTimeDisplay = ptime.New(p.StartTime).Format("HH:mm - yyyy/MM/dd")

	return p
}

func filterAttributes(attributes []schema.Taxonomy) (attrs []tresponse.Taxonomy) {
	for _, attr := range attributes {
		attrs = append(attrs, tresponse.Taxonomy{
			ID:       attr.ID,
			Title:    attr.Title,
			ParentID: attr.ParentID,
		})
	}
	return attrs
}
