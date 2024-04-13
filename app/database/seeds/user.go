package seeds

import (
	"errors"
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/helpers"
	"gorm.io/gorm"
)

type User struct{}

const UserSeedCount = 30
const AdminMobile = 9380338494

func (User) Seed(db *gorm.DB) error {
	pass := helpers.Hash([]byte("123456"))

	for i := 0; i <= UserSeedCount; i++ {
		fakeData := &schema.User{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.Password = pass
		fakeData.Permissions = schema.UserPermissions{}
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
	if count <= 1 {
		return 0, nil
	}

	return int(count), nil
}

func GenerateAdmin(db *gorm.DB) error {
	err := db.First(&schema.User{}, "mobile = ?", AdminMobile).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		return nil
	}

	pass := helpers.Hash([]byte("123456"))

	// create admin
	admin := &schema.User{}

	admin.Password = pass
	admin.LastName = "admin"
	admin.FirstName = "mahdi"
	admin.Mobile = AdminMobile
	admin.MobileConfirmed = true
	admin.Permissions = schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}

	return db.Create(admin).Error
}
