package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/helpers"
	"gorm.io/gorm"
)

type User struct{}

const UserSeedCount = 30

func (User) Seed(db *gorm.DB) error {
	pass := helpers.Hash([]byte("123456"))

	// create admin
	admin := &schema.User{}

	admin.Password = pass
	admin.LastName = "admin"
	admin.FirstName = "mahdi"
	admin.Mobile = 9380338494
	admin.MobileConfirmed = true
	admin.Roles = []string{schema.RAdmin, schema.RBusinessOwner}

	if err := db.Create(admin).Error; err != nil {
		log.Error().Err(err)
	}
	// end create admin

	for i := 0; i <= UserSeedCount; i++ {
		fakeData := &schema.User{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.Password = pass
		fakeData.Roles = []string{"user"}
		fakeData.Mobile = uint64(9180338500 + i)
		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d users created", UserSeedCount)

	return nil
}

func (User) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
