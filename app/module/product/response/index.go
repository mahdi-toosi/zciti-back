package response

import (
	"go-fiber-starter/app/database/schema"
	bresponse "go-fiber-starter/app/module/business/response"
	presponse "go-fiber-starter/app/module/post/response"
	tresponse "go-fiber-starter/app/module/taxonomy/response"
)

type Product struct {
	Post     presponse.Post
	Product  ProductInPost
	Variants []ProductInPost
}

type ProductInPost struct {
	ID          uint64
	Type        schema.ProductType        `json:",omitempty"`
	Price       float64                   `json:",omitempty"`
	OnSale      bool                      `json:",omitempty"`
	StockStatus schema.ProductStockStatus `json:",omitempty"`
	Meta        schema.ProductMeta        `json:",omitempty"`
	Attributes  []schema.ProductAttribute `json:",omitempty"`
}

func FromDomain(item *schema.Post, products []schema.Product) (res *Product) {
	if item == nil {
		return res
	}

	p := &Product{
		Post: presponse.Post{
			ID:       item.ID,
			Type:     item.Type,
			Meta:     item.Meta,
			Title:    item.Title,
			Status:   item.Status,
			Content:  item.Content,
			Business: bresponse.Business{ID: item.Business.ID, Title: item.Business.Title},
		},
		Variants: []ProductInPost{},
	}

	for _, taxonomy := range item.Taxonomies {
		p.Post.Taxonomies = append(p.Post.Taxonomies, tresponse.Taxonomy{
			ID:       taxonomy.ID,
			Type:     taxonomy.Type,
			Title:    taxonomy.Title,
			ParentID: taxonomy.ParentID,
		})
	}

	for _, product := range products {
		if product.IsRoot {
			p.Product = ProductInPost{
				ID:          product.ID,
				Type:        product.Type,
				Price:       product.Price,
				OnSale:      product.OnSale,
				StockStatus: product.StockStatus,
				Meta:        product.Meta,
				//Attributes:  product.Attributes,
			}
			continue
		}

		if product.Type != schema.ProductTypeVariant {
			continue
		}
		p.Variants = append(p.Variants, ProductInPost{
			ID:          product.ID,
			Type:        product.Type,
			Price:       product.Price,
			OnSale:      product.OnSale,
			StockStatus: product.StockStatus,
			Meta:        product.Meta,
			//Attributes:  product.Attributes,
		})
	}

	return p
}
