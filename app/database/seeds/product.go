package seeds

import (
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
)

type Product struct{}

const ProductSeedCount = 60

func (Product) Seed(db *gorm.DB) error {
	businessIDs, err := utils.GetFakeTableIDs(db, schema.Business{})
	if err != nil {
		return err
	}

	postIDs, err := utils.GetFakeTableIDsWithConditions(db, schema.Post{}, map[string][]any{"type": {schema.PostTypeProduct}}) // "productVariant"
	if err != nil {
		return err
	}

	for _, postID := range postIDs {
		fakeData := &schema.Product{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.IsRoot = true
		fakeData.Type = schema.ProductTypeVariant
		fakeData.VariantType = nil
		fakeData.PostID = postID
		fakeData.BusinessID = utils.RandomFromArray(businessIDs)

		if err := db.Create(&fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	for _, postID := range postIDs {
		fakeData := &schema.Product{}
		err := faker.FakeData(&fakeData)
		if err != nil {
			log.Error().Err(err).Msg("fail to generate fake data")
			return err
		}

		fakeData.PostID = postID
		fakeData.Type = schema.ProductTypeVariant
		fakeData.BusinessID = utils.RandomFromArray(businessIDs)

		if err := db.Create(&fakeData).Error; err != nil {
			log.Error().Err(err)
		}
	}

	log.Info().Msgf("%d products created", ProductSeedCount*2)

	return nil
}

func (Product) Count(db *gorm.DB) (int, error) {
	var count int64
	if err := db.Model(schema.Product{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
