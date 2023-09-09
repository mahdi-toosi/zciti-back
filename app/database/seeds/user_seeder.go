package seeds

import (
	"github.com/bangadam/go-fiber-starter/app/database/schema"
	"gorm.io/gorm"
)

type UserSeeder struct{}

var users = []schema.User{
	{
		FirstName:       "FirstName",
		LastName:        "LastName",
		Mobile:          "09999999999",
		MobileConfirmed: *bool(false),
		Password:        "123456",
	},
	{
		FirstName:       "FirstName",
		LastName:        "LastName",
		Mobile:          "09999999999",
		MobileConfirmed: *bool(false),
		Password:        "123456",
	},
}

func (UserSeeder) Seed(conn *gorm.DB) error {
	for _, row := range users {
		if err := conn.Create(&row).Error; err != nil {
			return err
		}
	}

	return nil
}

func (UserSeeder) Count(conn *gorm.DB) (int, error) {
	var count int64
	if err := conn.Model(&schema.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
