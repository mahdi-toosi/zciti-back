package schema

type Business struct {
	ID          uint64       `gorm:"primary_key" faker:"-"`
	Title       string       `gorm:"not null;varchar(255)" faker:"word"`
	Type        BusinessType `gorm:"not null;varchar(255)" faker:"oneof:GymManager,Bakery"`
	Description string       `gorm:"varchar(500)" faker:"paragraph"`
	OwnerID     uint64       `gorm:"not null" faker:"-"`
	Owner       User         `gorm:"foreignKey:OwnerID" faker:"-"`
	Users       []User       `gorm:"many2many:business_users;" faker:"-"`
	Base
}

type BusinessType string

const (
	TypeGymManager BusinessType = "GymManager"
	TypeBakery     BusinessType = "Bakery"
)

var TypeDisplayProxy = map[BusinessType]string{
	TypeGymManager: "مدیر باشگاه",
	TypeBakery:     "نانوایی",
}
