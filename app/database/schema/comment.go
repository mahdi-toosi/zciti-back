package schema

type Comment struct {
	ID              uint64        `gorm:"primaryKey" faker:"-"`
	Content         string        `gorm:"not null" faker:"paragraph"`
	Status          CommentStatus `gorm:"varchar(40); default:pending" faker:"oneof: approved, pending"`
	ParentID        uint64        `faker:"-"`
	AuthorID        uint64        `gorm:"not null" faker:"-"`
	Author          User          `gorm:"foreignKey:AuthorID" faker:"-"`
	IsBusinessOwner bool          `gorm:"not null" faker:"-"`
	PostID          uint64        `gorm:"index" faker:"-"`
	Post            Post          `gorm:"foreignKey:PostID" faker:"-"`
	Base
}

type CommentStatus string

const (
	CommentStatusPending  CommentStatus = "pending"
	CommentStatusApproved CommentStatus = "approved"
)
