package seeds

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
)

type Business struct{}

const BusinessSeedCount = 30

func (Business) Seed(db *gorm.DB) error {
	userIDs, err := utils.GetFakeTableIDs(db, schema.User{})
	if err != nil {
		return err
	}

	for i := 0; i <= BusinessSeedCount; i++ {
		fakeData := &schema.Business{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}
		fakeData.OwnerID = utils.RandomFromArray(userIDs)

		if err := db.Create(fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	businessIDs, err := utils.GetFakeTableIDs(db, schema.Business{})
	if err != nil {
		return err
	}

	query := "INSERT INTO business_users (business_id, user_id) VALUES  (%d, %d)"
	for b := 0; b < len(businessIDs); b++ {
		for u := 0; u < len(userIDs); u++ {
			db.Exec(fmt.Sprintf(query, businessIDs[b], userIDs[u]))
		}
	}

	log.Info().Msgf("%d businesses created", UserSeedCount)

	return nil
}

func (Business) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Business{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
