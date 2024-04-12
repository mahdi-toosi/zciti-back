package schema

type Post struct {
	ID             uint64   `gorm:"primaryKey" faker:"-"`
	Title          string   `gorm:"varchar(600);" faker:"sentence"`
	Content        string   `gorm:"not null" faker:"paragraph"`
	Status         string   `gorm:"varchar(100); default:published" faker:"oneof: draft, published, private"`
	Type           string   `gorm:"varchar(100);" faker:"oneof: product, post, page"`
	AuthorID       uint64   `gorm:"not null" faker:"-"`
	Author         User     `gorm:"foreignKey:AuthorID" faker:"-"`
	BusinessID     uint64   `faker:"-"`
	Business       Business `gorm:"foreignKey:BusinessID" faker:"-"`
	CommentsStatus string   `gorm:"varchar(100);" faker:"oneof: open, close, onlyBuyers, onlyCustomers"`
	CommentsCount  uint64   `gorm:"not null"`
	Base
}

const (
	PostStatusDraft     = "draft"
	PostStatusPublished = "published"
	PostStatusPrivate   = "private"
)

const (
	PostCommentStatusOpen              = "open"
	PostCommentStatusClose             = "close"
	PostCommentStatusOnlyBuyers        = "onlyBuyers"
	PostCommentStatusOnlyBusinessUsers = "onlyCustomers"
)

const (
	PostTypePage    = "page"
	PostTypeProduct = "product"
	PostTypePost    = "post"
)
