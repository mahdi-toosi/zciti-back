package seeds

import (
	"gorm.io/gorm"
)

type Seeder interface {
	Seed(*gorm.DB) error
	Count(*gorm.DB) (int, error)
}

func MainDBSeeders() []Seeder {
	// order matters
	return []Seeder{
		User{},
		Business{},
		Post{},
		Comment{},
		NotificationTemplate{},
		Notification{},
	}
}

func ChatDBSeeders() []Seeder {
	// order matters
	return []Seeder{
		MessageRoom{},
		Message{},
	}
}

func GenerateNecessaryData(db *gorm.DB) error {
	if err := GenerateAdmin(db); err != nil {
		return err
	}

	if err := GenerateRootBusiness(db); err != nil {
		return err
	}

	return nil
}
