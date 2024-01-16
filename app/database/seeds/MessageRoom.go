package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
)

type MessageRoom struct{}

const MessageRoomSeedCount = 30

func (MessageRoom) Seed(db *gorm.DB) error {
	for i := 0; i <= MessageRoomSeedCount; i++ {
		fakeData := &schema.MessageRoom{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.UserID = 1 // admin
		fakeData.BusinessID = utils.Random(1, BusinessSeedCount)

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d message rooms created", MessageRoomSeedCount)

	return nil
}

func (MessageRoom) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.MessageRoom{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
