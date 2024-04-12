package schema

type Business struct {
	ID          uint64          `gorm:"primaryKey" faker:"-"`
	Title       string          `gorm:"not null;varchar(255)" faker:"word"`
	Type        BusinessType    `gorm:"not null;varchar(255)" faker:"oneof:GymManager,Bakery"`
	Description string          `gorm:"varchar(500)" faker:"paragraph"`
	OwnerID     uint64          `gorm:"not null" faker:"-"`
	Owner       User            `gorm:"foreignKey:OwnerID" faker:"-"`
	Users       []*User         `gorm:"many2many:business_users;" faker:"-"`
	Account     BusinessAccount `gorm:"varchar(100);default:default" faker:"-"`
	AssetsSize  uint64          `gorm:"default:0" faker:"-"`
	ShebaNumber string          `gorm:"varchar(24)" faker:"-"`
	Base
}

type BusinessType string

const (
	BTypeROOT       BusinessType = "ROOT"
	BTypeBakery     BusinessType = "Bakery"
	BTypeGymManager BusinessType = "GymManager"
)

var TypeDisplayProxy = map[BusinessType]string{
	BTypeGymManager: "مدیر باشگاه",
	BTypeBakery:     "نانوایی",
}

type BusinessAccount string

const (
	BusinessAccountDefault BusinessAccount = "default"
)

const ROOT_BUSINESS_ID = 1
