package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
)

type Post struct{}

const PostSeedCount = 60

func (Post) Seed(db *gorm.DB) error {
	for i := 0; i <= PostSeedCount; i++ {
		fakeData := &schema.Post{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.AuthorID = utils.Random(1, UserSeedCount)
		fakeData.BusinessID = utils.Random(1, BusinessSeedCount)

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d posts created", PostSeedCount)

	return nil
}

func (Post) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Post{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
