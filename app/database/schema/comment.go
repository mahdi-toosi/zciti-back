package schema

type Comment struct {
	ID       uint64 `gorm:"primaryKey" faker:"-"`
	Content  string `gorm:"not null" faker:"paragraph"`
	Status   string `gorm:"varchar(40); default:pending" faker:"oneof: approved, pending"`
	AuthorID uint64 `gorm:"not null" faker:"-"`
	Author   User   `gorm:"foreignKey:AuthorID" faker:"-"`
	PostID   uint64 `faker:"-"`
	Base
}
