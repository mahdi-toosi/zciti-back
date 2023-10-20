package schema

type Post struct {
	ID       uint64 `gorm:"primary_key" faker:"-"`
	Title    string `gorm:"varchar(600);" faker:"sentence"`
	Content  string `gorm:"not null" faker:"paragraph"`
	Status   string `gorm:"varchar(100); default:published" faker:"oneof: draft, published, private"`
	Type     string `gorm:"varchar(100);" faker:"oneof: product, article, page"`
	AuthorID uint64 `gorm:"not null" faker:"-"`
	Author   User   `gorm:"foreignKey:AuthorID" faker:"-"`
	Base
}

// BusinessID uint64 `gorm:"not null" faker:"-"`
// Business Business `gorm:"foreignKey:BusinessID" faker:"-"`
