package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
)

type Order struct{}

const OrderSeedCount = 100

func (Order) Seed(db *gorm.DB) error {
	businessIDs, err := utils.GetFakeTableIDs(db, schema.Business{})
	if err != nil {
		return err
	}
	userIDs, err := utils.GetFakeTableIDs(db, schema.User{})
	if err != nil {
		return err
	}

	for i := 0; i <= OrderSeedCount; i++ {
		fakeData := &schema.Order{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.UserID = utils.RandomFromArray(userIDs)
		fakeData.BusinessID = utils.RandomFromArray(businessIDs)

		if err := db.Create(&fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d orders created", OrderSeedCount)

	return nil
}

func (Order) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Order{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
