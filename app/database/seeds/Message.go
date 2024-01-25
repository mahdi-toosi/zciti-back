package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
)

type Message struct{}

const MessageSeedCount = 400

func (Message) Seed(db *gorm.DB) error {
	var rooms []schema.MessageRoom
	err := db.Model(&schema.MessageRoom{}).Limit(5).Find(&rooms).Error
	if err != nil {
		return err
	}

	for i := 0; i <= MessageSeedCount; i++ {
		fakeData := &schema.Message{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.SeenAt = utils.RandomDateTime()
		fakeData.FromID = 1 // admin
		fakeData.RoomID = rooms[int(utils.Random(0, 3))].ID

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d message created", MessageSeedCount)

	return nil
}

func (Message) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Message{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
