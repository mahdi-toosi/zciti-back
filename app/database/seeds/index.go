package seeds

import "gorm.io/gorm"

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
