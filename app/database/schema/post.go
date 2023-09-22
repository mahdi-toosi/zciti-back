package schema

type Post struct {
	ID       uint64 `gorm:"primary_key" faker:"-"`
	AuthorID uint64 `gorm:"not null" faker:"oneof: 1, 2 ,3 ,4 ,5 ,6 ,7 ,8 ,9 ,10"`
	Title    string `gorm:"varchar(600);" faker:"sentence"`
	Content  string `faker:"paragraph"`
	Status   string `gorm:"varchar(100); default:published" faker:"oneof: draft, published, private"`
	Type     string `gorm:"varchar(100);" faker:"oneof: product, article, page"`
	Base
}

// BusinessID uint64 `faker:"oneof: 1, 2 ,3 ,4 ,5 ,6 ,7 ,8 ,9 ,10"`
