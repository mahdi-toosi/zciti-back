package response

import (
	"go-fiber-starter/app/database/schema"
	tresponse "go-fiber-starter/app/module/taxonomy/response"
	"time"
)

type Reservation struct {
	ID            uint64    `json:",omitempty"`
	EndTime       time.Time `json:",omitempty"`
	StartTime     time.Time `json:",omitempty"`
	ProductID     uint64    `json:",omitempty"`
	ProductSKU    string    `json:",omitempty"`
	ProductDetail string    `json:",omitempty"`
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
		ProductDetail: item.Product.Meta.Detail,
	}

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
