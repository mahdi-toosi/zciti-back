package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Product struct {
	ID          uint64              `gorm:"primaryKey" faker:"-"`
	PostID      uint64              `gorm:"index" faker:"-"`
	Post        Post                `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" faker:"-"`
	IsRoot      bool                `faker:"-"`
	Type        ProductType         `gorm:"varchar(50); not null" faker:"oneof: variant"` // simple, variant ,grouped, reservable, downloadable
	VariantType *ProductVariantType `gorm:"varchar(50);" faker:"oneof: washingMachine"`   //  simple, reservable, downloadable
	Price       float64             `gorm:"not null"`                                     // for variants
	MinPrice    float64             `gorm:"not null"`                                     // for variants
	MaxPrice    float64             `gorm:"not null"`                                     // for variants
	OnSale      bool                ``
	StockStatus ProductStockStatus  `gorm:"varchar(40); not null;" faker:"oneof: inStock, outOfStock, onBackOrder"`
	TotalSales  float64             ``
	Meta        ProductMeta         `gorm:"type:jsonb"`
	Taxonomies  []Taxonomy          `gorm:"many2many:products_taxonomies;" faker:"-"` // `gorm:"foreignKey:ProductID"`
	BusinessID  uint64              `gorm:"index" faker:"-"`
	Business    Business            `gorm:"foreignKey:BusinessID" faker:"-"`
	//Attributes  []ProductAttribute  `gorm:"many2many:products_attributes;" faker:"-"` // `gorm:"foreignKey:ProductID"`
	Base
}

type ProductType string

const (
	ProductTypeSimple       ProductType = "simple"     // A simple product is a standalone physical or digital product that doesn't have any variations. For example, a book or a music download.
	ProductTypeVariant      ProductType = "variant"    // A variable product is a product that has multiple variations, such as different sizes or colors. Each variation can have its own price, SKU, and stock level. For example, a t-shirt that is available in different sizes and colors.
	ProductTypeGrouped      ProductType = "grouped"    // A grouped product is a collection of related products that are sold together. For example, a computer bundle that includes a monitor, keyboard, and mouse.
	ProductTypeVirtual      ProductType = "virtual"    // A virtual product is a non-physical product that doesn't require shipping. For example, a consulting service or an online course.
	ProductTypeReservable   ProductType = "reservable" // A reservable product is a product that can be reserved by customers. For example, a hotel room that can be booked for a specific date and time.
	ProductTypeDownloadable ProductType = "downloadable"
)

type ProductVariantType string

const (
	ProductVariantTypeSimple         ProductVariantType = "simple"
	ProductVariantTypeReservable     ProductVariantType = "reservable"
	ProductVariantTypeDownloadable   ProductVariantType = "downloadable"
	ProductVariantTypeWashingMachine ProductVariantType = "washingMachine"
)

type ProductStockStatus string

const (
	ProductStockStatusInStock     ProductStockStatus = "inStock"
	ProductStockStatusOutOfStock  ProductStockStatus = "outOfStock"
	ProductStockStatusOnBackorder ProductStockStatus = "onBackOrder" // In WooCommerce, the "on backorder" stock status indicates that a product is currently out of stock but more stock is expected to arrive at a later date. Customers can still place orders for the product while it is on backorder, and the order will be fulfilled when the new stock arrives.
)

type ProductTaxStatus string

const (
	ProductTaxStatusNone     ProductTaxStatus = "none"
	ProductTaxStatusTaxable  ProductTaxStatus = "taxable"
	ProductTaxStatusShipping ProductTaxStatus = "shipping" // In WooCommerce, the "shipping" tax status indicates that the product is taxable but the tax rate is calculated based on the shipping cost.
)

type ProductMetaSelectedAttribute struct {
	ID                uint64
	Title             string
	ShowInProductPage bool
	UseInVariants     bool
}

type ProductMetaReservationInfoData struct {
	To   string // "13:22:00"
	From string // "13:22:00"
}

type ProductMetaReservation struct {
	Quantity int           `json:",omitempty" example:"1"`
	Duration time.Duration `json:",omitempty" example:"1"`

	EndTime   string `json:",omitempty" example:"1"` // "13:22:00"
	StartTime string `json:",omitempty" example:"1"` // "13:22:00"

	Info map[time.Weekday] /* week day num*/ []ProductMetaReservationInfoData `json:",omitempty" example:"1"`
}

type UniWashMachineStatus string

const (
	UniWashMachineStatusON  UniWashMachineStatus = "ON"
	UniWashMachineStatusOFF UniWashMachineStatus = "OFF"
)

type ProductMeta struct {
	Detail               string                         `json:",omitempty" example:"Detail" validate:"omitempty,min=2,max=500" faker:"word"`  // The stock keeping unit (SKU) of the product. This is a unique identifier for the product that is used for inventory management.
	SKU                  string                         `json:",omitempty" example:"sku-2f3s" validate:"omitempty,min=2,max=40" faker:"word"` // The stock keeping unit (SKU) of the product. This is a unique identifier for the product that is used for inventory management.
	UniWashMobileNumber  string                         `json:",omitempty" example:"09909999999"  faker:"-"`                                  //validate:"omitempty,min=10,max=10"                                // The stock keeping unit (SKU) of the product. This is a unique identifier for the product that is used for inventory management.
	UniWashMachineStatus UniWashMachineStatus           `json:",omitempty" faker:"-"`
	PurchaseNote         string                         `json:",omitempty" validator:"omitempty,min=2,max=500"` // A note that is displayed to the customer after purchasing the product.
	Weight               float64                        `json:",omitempty" validator:"omitempty,number"`
	Width                float64                        `json:",omitempty" validator:"omitempty,number"`
	Height               float64                        `json:",omitempty" validator:"omitempty,number"`
	Length               float64                        `json:",omitempty" validator:"omitempty,number"`
	ProviderPrice        float64                        `json:",omitempty" validator:"omitempty,number"`
	CouldReserveUntil    time.Time                      `json:",omitempty" validator:"omitempty,datetime"` // millisecond from now
	TaxStatus            ProductTaxStatus               `json:",omitempty" validator:"omitempty,oneof: none taxable shipping" faker:"oneof: none, taxable, shipping"`
	Images               []string                       `json:",omitempty" faker:"-"`
	AttributesMap        map[uint64]uint64              `faker:"-"`
	SelectedAttributes   []ProductMetaSelectedAttribute `json:",omitempty" faker:"-"`

	SalePrice          float64   `json:",omitempty" validator:"omitempty,number"`
	SalePriceStartDate time.Time `json:",omitempty" validator:"omitempty,datetime"`
	SalePriceEndDate   time.Time `json:",omitempty" validator:"omitempty,datetime"`

	ManageStock   bool   `json:",omitempty"`
	StockSku      string `json:",omitempty" example:"sku-2f3s" validate:"omitempty,min=2,max=40" faker:"word"`
	StockQuantity uint64 `json:",omitempty" validate:"omitempty,number"` // The number of units of the product that are currently in stock.

	Reservation ProductMetaReservation `json:",omitempty" faker:"-"`
}

func (pm *ProductMeta) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal ProductMeta with value %v", value)
	}
	return json.Unmarshal(byteValue, pm)
}

func (pm ProductMeta) Value() (driver.Value, error) {
	return json.Marshal(pm)
}
