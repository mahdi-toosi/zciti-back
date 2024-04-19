package seeds

import (
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"gorm.io/gorm"
	"time"
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

	time.Sleep(time.Second * 1)
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
	if count <= 1 {
		return 0, nil
	}

	return int(count), nil
}

func GenerateRootBusiness(db *gorm.DB) error {
	const title = "Root"
	err := db.First(&schema.Business{}, "title = ? AND type = ?", title, schema.BTypeROOT).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		return nil
	}

	business := schema.Business{OwnerID: 1, Title: title, Type: schema.BTypeROOT}
	if err = db.Create(&business).Error; err != nil {
		return err
	}

	admin := schema.User{}
	if err = db.First(&admin, "mobile = ?", AdminMobile).Error; err != nil {
		return err
	}

	query := "INSERT INTO business_users (business_id, user_id) VALUES  (%d, %d)"
	return db.Exec(fmt.Sprintf(query, business.ID, admin.ID)).Error
}
