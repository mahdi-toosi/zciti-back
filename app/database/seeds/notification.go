package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
)

type Notification struct{}

const NotificationSeedCount = 40

func (Notification) Seed(db *gorm.DB) error {
	for i := 0; i <= NotificationSeedCount; i++ {
		fakeData := &schema.Notification{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		if i%2 == 0 {
			fakeData.Type = []string{string(schema.TNotification)}
		} else {
			fakeData.Type = []string{string(schema.TSms)}
		}

		fakeData.ReceiverID = utils.Random(1, UserSeedCount)
		fakeData.TemplateID = utils.Random(1, NotificationTemplateSeedCount)

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d notifications created", NotificationSeedCount)

	return nil
}

func (Notification) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Notification{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
