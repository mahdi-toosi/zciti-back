package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/utils/paginator"
)

type Product struct {
	Post     request.Post
	Product  ProductInPost
	Variants []ProductInPost
}

type ProductInPost struct {
	ID          uint64
	PostID      uint64                    `example:"1" validate:"number"`
	Type        schema.ProductType        `example:"simple" validate:"required,oneof=simple variant washingMachine"`
	Price       float64                   `example:"65000" validate:"required,number"`
	OnSale      bool                      `example:"true"`
	StockStatus schema.ProductStockStatus `example:"inStock" validate:"required,oneof=inStock outOfStock onBackOrder"`
	Meta        schema.ProductMeta
	Attributes  []schema.ProductAttribute
}

type ProductsRequest struct {
	BusinessID uint64
	Keyword    string
	Pagination *paginator.Pagination
}

func (req *Product) ToDomain(postID uint64, businessID uint64) (products []*schema.Product) {
	var MinPrice = req.Product.Price
	var MaxPrice = req.Product.Price

	if req.Product.Type == schema.ProductTypeVariant {
		for _, product := range req.Variants {
			if product.Price < MinPrice {
				MinPrice = product.Price
			}
			if product.Price > MaxPrice {
				MaxPrice = product.Price
			}

			products = append(products, &schema.Product{
				PostID:     postID,
				MinPrice:   MinPrice,
				MaxPrice:   MaxPrice,
				BusinessID: businessID,
				ID:         product.ID,
				Meta:       product.Meta,
				Type:       product.Type,
				OnSale:     product.OnSale,
				//Attributes:  product.Attributes,
				StockStatus: product.StockStatus,
			})
		}
	}

	products = append(products, &schema.Product{
		IsRoot:     true,
		PostID:     postID,
		MinPrice:   MinPrice,
		MaxPrice:   MaxPrice,
		BusinessID: businessID,
		ID:         req.Product.ID,
		Price:      req.Product.Price,
		Meta:       req.Product.Meta,
		Type:       req.Product.Type,
		OnSale:     req.Product.OnSale,
		//Attributes:  req.Product.Attributes,
		StockStatus: req.Product.StockStatus,
	})

	return products
}
