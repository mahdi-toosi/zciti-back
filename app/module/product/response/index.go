package response

import (
	"go-fiber-starter/app/database/schema"
	presponse "go-fiber-starter/app/module/post/response"
	tresponse "go-fiber-starter/app/module/taxonomy/response"
)

type Product struct {
	Product  ProductInPost
	Post     presponse.Post
	Variants []ProductInPost
}

type ProductInPost struct {
	ID          uint64
	Price       float64                    `json:",omitempty"`
	OnSale      bool                       `json:",omitempty"`
	Type        schema.ProductType         `json:",omitempty"`
	Meta        schema.ProductMeta         `json:",omitempty"`
	Taxonomies  []tresponse.Taxonomy       `json:",omitempty"`
	StockStatus schema.ProductStockStatus  `json:",omitempty"`
	VariantType *schema.ProductVariantType `json:",omitempty"`
	Attributes  []tresponse.Taxonomy       `json:",omitempty"`
}

func FromDomain(item *schema.Post, products []schema.Product, isForUser bool) (res *Product) {
	if item == nil {
		return res
	}

	if isForUser {
		p := &Product{
			Post: presponse.Post{
				ID:      item.ID,
				Meta:    item.Meta,
				Title:   item.Title,
				Excerpt: item.Excerpt,
				Content: item.Content,
				//Business: bresponse.Business{ID: item.Business.ID, Title: item.Business.Title},
			},
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
					Meta:        product.Meta,
					Price:       product.Price,
					OnSale:      product.OnSale,
					StockStatus: product.StockStatus,
					//Attributes:  filterAttributes(product.Taxonomies),
				}
				continue
			}

			g := ProductInPost{
				ID:          product.ID,
				Type:        product.Type,
				Meta:        product.Meta,
				Price:       product.Price,
				OnSale:      product.OnSale,
				StockStatus: product.StockStatus,
				VariantType: product.VariantType,
				//Attributes:  filterAttributes(product.Taxonomies),
			}
			p.Variants = append(p.Variants, g)
		}

		return p
	}

	p := &Product{
		Post: presponse.Post{
			ID:      item.ID,
			Type:    item.Type,
			Meta:    item.Meta,
			Title:   item.Title,
			Status:  item.Status,
			Excerpt: item.Excerpt,
			Content: item.Content,
			//Business: bresponse.Business{ID: item.Business.ID, Title: item.Business.Title},
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
				Meta:        product.Meta,
				Price:       product.Price,
				OnSale:      product.OnSale,
				StockStatus: product.StockStatus,
				//Attributes:  filterAttributes(product.Taxonomies),
			}
			continue
		}

		p.Variants = append(p.Variants, ProductInPost{
			ID:          product.ID,
			Type:        product.Type,
			Meta:        product.Meta,
			Price:       product.Price,
			OnSale:      product.OnSale,
			VariantType: product.VariantType,
			StockStatus: product.StockStatus,
			//Attributes:  filterAttributes(product.Taxonomies),
		})
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
