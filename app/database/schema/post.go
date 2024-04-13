package schema

type Post struct {
	ID             uint64            `gorm:"primaryKey" faker:"-"`
	Title          string            `gorm:"varchar(600);" faker:"sentence"`
	Content        string            `gorm:"not null" faker:"paragraph"`
	Status         PostStatus        `gorm:"varchar(100); default:published" faker:"oneof: draft, published, private"`
	Type           PostType          `gorm:"varchar(100);" faker:"oneof: product, post, page"`
	AuthorID       uint64            `gorm:"not null" faker:"-"`
	Author         User              `gorm:"foreignKey:AuthorID" faker:"-"`
	BusinessID     uint64            `faker:"-"`
	Business       Business          `gorm:"foreignKey:BusinessID" faker:"-"`
	Taxonomies     []Taxonomy        `gorm:"many2many:post_taxonomy;" faker:"-"`
	CommentsStatus PostCommentStatus `gorm:"varchar(100);" faker:"oneof: open, close, onlyBuyers, onlyCustomers"`
	CommentsCount  uint64            `gorm:"not null"`
	Base
}

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusPrivate   PostStatus = "private"
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
	PostTypeProduct PostType = "product"
	PostTypePost    PostType = "post"
)
