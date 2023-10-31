package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"gorm.io/gorm"
)

type NotificationTemplate struct{}

const NotificationTemplateSeedCount = 30

func (NotificationTemplate) Seed(db *gorm.DB) error {
	for i := 0; i <= NotificationTemplateSeedCount; i++ {
		fakeData := &schema.NotificationTemplate{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}
	faker.ResetUnique()

	log.Info().Msgf("%d notification templates created", NotificationTemplateSeedCount)

	return nil
}

func (NotificationTemplate) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.NotificationTemplate{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
