package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"gorm.io/gorm"
)

type Article struct{}

func (Article) Seed(conn *gorm.DB) error {
	count := 30
	for i := 0; i <= count; i++ {
		fakeData := &schema.Post{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		if err := conn.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d users created", count)

	return nil
}

func (Article) Count(conn *gorm.DB) (int, error) {
	var count int64
	if err := conn.Model(schema.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
