package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/helpers"
	"gorm.io/gorm"
)

type User struct{}

func (User) Seed(db *gorm.DB) error {
	count := 30
	for i := 0; i <= count; i++ {
		fakeData := &schema.User{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.Mobile = uint64(9180338500 + i)
		fakeData.Roles = []string{"user"}
		fakeData.Password = helpers.Hash([]byte("123456"))
		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}
	faker.ResetUnique()

	// create admin
	admin := &schema.User{}
	admin.LastName = "admin"
	admin.FirstName = "mahdi"
	admin.Mobile = 9180338595
	admin.MobileConfirmed = true
	admin.Roles = []string{"admin"}
	admin.Password = helpers.Hash([]byte("123456"))

	if err := db.Create(admin).Error; err != nil {
		log.Error().Err(err)
	}
	// end create admin

	log.Info().Msgf("%d users created", count)

	return nil
}

func (User) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
