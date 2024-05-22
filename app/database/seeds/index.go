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
		Product{},
		//Order{},
		Comment{},
		Taxonomy{},
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

	if err := GenerateUniWashOperator(db); err != nil {
		return err
	}

	if err := GenerateRootBusiness(db); err != nil {
		return err
	}

	if err := GenerateUniWashBusiness(db); err != nil {
		return err
	}

	return nil
}
