package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"gorm.io/gorm"
)

type Post struct{}

func (Post) Seed(db *gorm.DB) error {
	log.Info().Msg("mahdi yahooooo")
	count := 30
	for i := 0; i <= count; i++ {
		fakeData := &schema.Post{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d users created", count)

	return nil
}

func (Post) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Post{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
