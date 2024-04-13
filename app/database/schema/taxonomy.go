package schema

type Taxonomy struct {
	ID         uint64       `gorm:"primaryKey" faker:"-"`
	Name       string       `gorm:"varchar(100); not null;" faker:"word"`
	Type       TaxonomyType `gorm:"varchar(100); not null;" faker:"word"`
	IsGeneral  bool         ``
	Slug       string       `gorm:"varchar(200); not null;" faker:"word"`
	BusinessID uint64       `faker:"-"`
	Business   Business     `gorm:"foreignKey:BusinessID" faker:"-"`
	Posts      []Post       `gorm:"many2many:post_taxonomy;" faker:"-"`
	ParentID   uint64       `faker:"-"`
	Base
}

type TaxonomyType string

const (
	Category TaxonomyType = "category"
	Tag      TaxonomyType = "tag"
)
