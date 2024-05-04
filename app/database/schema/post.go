package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
)

type Post struct {
	ID         uint64     `gorm:"primaryKey" faker:"-"`
	Title      string     `gorm:"varchar(255);" faker:"sentence"`
	Excerpt    string     `gorm:"varchar(255);" faker:"sentence"`
	Content    string     `gorm:"not null" faker:"paragraph"`
	Status     PostStatus `gorm:"varchar(50); default:published" faker:"oneof: draft, published"`
	Type       PostType   `gorm:"varchar(50); not null;" faker:"oneof: post, page, product"`
	ParentID   uint64     `faker:"-"`
	Slug       string     `gorm:"varchar(600); index:idx_slug;" faker:"-"`
	AuthorID   uint64     `gorm:"not null" faker:"-"`
	Author     User       `gorm:"foreignKey:AuthorID" faker:"-"`
	BusinessID uint64     `gorm:"index:idx_slug;" faker:"-"`
	Business   Business   `gorm:"foreignKey:BusinessID" faker:"-"`
	Products   []Product  `gorm:"foreignKey:PostID" faker:"-"`
	Taxonomies []Taxonomy `gorm:"many2many:posts_taxonomies;" faker:"-"`
	Meta       PostMeta   `gorm:"type:jsonb"`
	Base
}

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	//PostStatusPrivate   PostStatus = "private"
)

type PostCommentStatus string

const (
	PostCommentStatusOpen              PostCommentStatus = "open"
	PostCommentStatusClose             PostCommentStatus = "close"
	PostCommentStatusOnlyBuyers        PostCommentStatus = "onlyBuyers"
	PostCommentStatusOnlyBusinessUsers PostCommentStatus = "onlyCustomers"
)

type PostType string

const (
	PostTypePage    PostType = "page"
	PostTypePost    PostType = "post"
	PostTypeProduct PostType = "product"
)

type PostMeta struct {
	CommentsStatus PostCommentStatus `example:"open" validate:"required,oneof=open close onlyBuyers onlyCustomers" faker:"oneof: open, close, onlyBuyers, onlyCustomers"`
	CommentsCount  uint64            // `gorm:"not null"`
	FeaturedImage  string            `faker:"-"`
	//---
	UpSellIDs    string `json:",omitempty"`                                     // 1,2,3 related products that is better, For example, if a customer is looking to buy a basic smartphone, an upsell might be to offer them a more advanced model with more features and a higher price point
	CrossSellIDs string `json:",omitempty"`                                     // 1,2,3 related products , For example, if a customer is buying a camera, a cross-sell might be to offer them a memory card, camera case, or other related accessories.
	PurchaseNote string `validator:"omitempty,min=2,max=500" json:",omitempty"` // A note that is displayed to the customer after purchasing the product.
}

func (pm *PostMeta) Scan(value any) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal PostMeta with value %v", value)
	}
	return json.Unmarshal(byteValue, pm)
}

func (pm PostMeta) Value() (driver.Value, error) {
	return json.Marshal(pm)
}

func (u *Post) GenerateSlug() string {
	return fmt.Sprintf("%s-%d", slug.Make(u.Title), u.BusinessID)
}
