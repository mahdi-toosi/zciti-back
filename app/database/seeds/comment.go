package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
)

type Comment struct{}

const CommentSeedCount = 400

func (Comment) Seed(db *gorm.DB) error {
	userIDs, err := utils.GetFakeTableIDs(db, schema.User{})
	if err != nil {
		return err
	}
	postIDs, err := utils.GetFakeTableIDs(db, schema.Post{})
	if err != nil {
		return err
	}

	for i := 0; i <= CommentSeedCount; i++ {
		fakeData := &schema.Comment{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.PostID = utils.RandomFromArray(postIDs)
		fakeData.AuthorID = utils.RandomFromArray(userIDs)

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d comments created", CommentSeedCount)

	return nil
}

func (Comment) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Comment{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
