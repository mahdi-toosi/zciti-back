package schema

type Post struct {
	ID              uint64 `gorm:"primary_key" faker:"-"`
	FirstName       string `faker:"first_name"`
	LastName        string `faker:"last_name"`
	Mobile          string `gorm:"not null;uniqueIndex" faker:"e_164_phone_number"`
	MobileConfirmed bool   `gorm:"default:false"`
	Password        string `gorm:"not null" faker:"password"`
	Base
}
