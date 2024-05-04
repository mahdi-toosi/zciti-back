package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/utils/paginator"
)

type StoreProductAttribute struct {
	ProductID     uint64 `example:"1" validate:"number"`
	BusinessID    uint64 `example:"1"`
	AddedAttrID   uint64 `example:"1" validate:"number"`
	RemovedAttrID uint64 `example:"1" validate:"omitempty,number"`
}

type Product struct {
	Post     request.Post
	Product  ProductInPost
	Variants []ProductInPost
}

type ProductInPost struct {
	ID          uint64
	OnSale      bool                       `example:"true"`
	PostID      uint64                     `example:"1" validate:"number"`
	BusinessID  uint64                     `example:"1"`
	Price       float64                    `example:"65000" validate:"omitempty,number"`
	Type        schema.ProductType         `example:"simple" validate:"required,oneof=simple variant"`
	StockStatus schema.ProductStockStatus  `example:"inStock" validate:"required,oneof=inStock outOfStock onBackOrder"`
	VariantType *schema.ProductVariantType `example:"simple" validate:"omitempty,oneof=simple reservable downloadable washingMachine"`
	Meta        schema.ProductMeta
}

type ProductsRequest struct {
	BusinessID uint64
	Keyword    string
	Pagination *paginator.Pagination
}

func (req *Product) ToDomain(postID uint64, businessID uint64) (products []*schema.Product) {
	var MinPrice = req.Product.Price
	var MaxPrice = req.Product.Price

	for _, variant := range req.Variants {

		if variant.Price <= MinPrice {
			MinPrice = variant.Price
		}
		if variant.Price >= MaxPrice {
			MaxPrice = variant.Price
		}

		products = append(products, &schema.Product{
			PostID:      postID,
			MinPrice:    MinPrice,
			MaxPrice:    MaxPrice,
			BusinessID:  businessID,
			ID:          variant.ID,
			Meta:        variant.Meta,
			Type:        variant.Type,
			Price:       variant.Price,
			OnSale:      variant.OnSale,
			StockStatus: variant.StockStatus,
			VariantType: variant.VariantType,
			//Attributes:  product.Attributes,
		})
	}

	products = append(products, &schema.Product{
		IsRoot:      true,
		PostID:      postID,
		MinPrice:    MinPrice,
		MaxPrice:    MaxPrice,
		BusinessID:  businessID,
		ID:          req.Product.ID,
		Price:       req.Product.Price,
		Meta:        req.Product.Meta,
		Type:        req.Product.Type,
		OnSale:      req.Product.OnSale,
		StockStatus: req.Product.StockStatus,
		//Attributes:  req.Product.Attributes,
	})

	return products
}

func (req *ProductInPost) ToDomain(postID uint64, businessID uint64) (product *schema.Product) {
	product = &schema.Product{
		PostID:      postID,
		ID:          req.ID,
		Meta:        req.Meta,
		Type:        req.Type,
		Price:       req.Price,
		OnSale:      req.OnSale,
		BusinessID:  businessID,
		StockStatus: req.StockStatus,
		VariantType: req.VariantType,
	}

	return product
}
