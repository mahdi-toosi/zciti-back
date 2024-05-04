package schema

import (
	"fmt"
	"github.com/gosimple/slug"
)

type Taxonomy struct {
	ID          uint64       `gorm:"primaryKey" faker:"-"`
	Title       string       `gorm:"varchar(100); not null;" faker:"word"`
	Type        TaxonomyType `gorm:"varchar(100); not null; index;" faker:"oneof: tag, category, productAttribute"`
	Domain      PostType     `gorm:"varchar(100); not null; index;" faker:"oneof: post, page, product"`
	Slug        string       `gorm:"varchar(200); not null; index:idx_slug,priority:2" faker:"-"`
	BusinessID  uint64       `gorm:"index:idx_slug,priority:1" faker:"-"`
	Business    Business     `gorm:"foreignKey:BusinessID" faker:"-"`
	Posts       []Post       `gorm:"many2many:posts_taxonomies;" faker:"-"`
	Products    []Product    `gorm:"many2many:products_taxonomies;" faker:"-"`
	ParentID    *uint64      `faker:"-"`
	Description string       `gorm:"varchar(500);" faker:"sentence"`
	Base
}

type TaxonomyType string

const (
	TaxonomyTypeTag               TaxonomyType = "tag"
	TaxonomyTypeCategory          TaxonomyType = "category"
	TaxonomyTypeProductAttributes TaxonomyType = "productAttribute"
)

func (u *Taxonomy) GenerateSlug() string {
	return fmt.Sprintf("%s-%d", slug.Make(u.Title), u.BusinessID)
}
