package seeds

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
	"time"
)

type Taxonomy struct{}

const TaxonomySeedCount = 100

func (Taxonomy) Seed(db *gorm.DB) error {
	businessIDs, err := utils.GetFakeTableIDs(db, schema.Business{})
	if err != nil {
		return err
	}

	for i := 0; i <= TaxonomySeedCount; i++ {
		fakeData := &schema.Taxonomy{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.BusinessID = utils.RandomFromArray(businessIDs)
		fakeData.Slug = fakeData.GenerateSlug() + "-" + utils.RandomStringBytes(5)

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	time.Sleep(time.Second * 1)

	postsIDs, err := utils.GetFakeTableIDs(db, schema.Post{})
	if err != nil {
		return err
	}
	taxonomiesIDs, err := utils.GetFakeTableIDs(db, schema.Taxonomy{})
	if err != nil {
		return err
	}

	query := "INSERT INTO post_taxonomy (taxonomy_id, post_id) VALUES  (%d, %d)"
	for u := 0; u < len(postsIDs); u++ {
		for b := 0; b < int(utils.Random(int(utils.Random(0, len(taxonomiesIDs))), len(taxonomiesIDs))); b++ {
			err = db.Exec(fmt.Sprintf(query, taxonomiesIDs[b], postsIDs[u])).Error
			if err != nil {
				log.Error().Err(err)
				return err
			}
		}
	}

	log.Info().Msgf("%d taxonomies created", TaxonomySeedCount)

	return nil
}

func (Taxonomy) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Taxonomy{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
