package seeder

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/utils/database/faker"
	"gorm.io/gorm"
)

type Seeder struct {
	Seeder interface{}
}

func RegisterSeeders(db *gorm.DB) []Seeder {
	return []Seeder{
		{Seeder: faker.UserFaker(db)},
		{Seeder: faker.ProductFaker(db)},
		{Seeder: faker.ProductPhotoFaker(db)},
		{Seeder: faker.CategoryFaker(db)},
		{Seeder: faker.CarouselFaker(db)},
		{Seeder: faker.ArticleFaker(db)},
		{Seeder: faker.ChallengeFaker(db)},
		{Seeder: faker.VoucherFaker(db)},
	}
}

func IsSeederExecuted(db *gorm.DB) bool {
	var status entities.StatusSeederModels
	result := db.First(&status)
	if result.Error != nil {
		return false
	}
	return status.IsExecuted
}

func SetSeederStatus(db *gorm.DB, executed bool) error {
	status := entities.StatusSeederModels{IsExecuted: executed}
	result := db.Create(&status)
	return result.Error
}

func DBSeed(db *gorm.DB) error {
	if !IsSeederExecuted(db) {
		for _, seeder := range RegisterSeeders(db) {
			err := db.Debug().Create(seeder.Seeder).Error
			if err != nil {
				return err
			}
		}

		err := SetSeederStatus(db, true)
		if err != nil {
			return err
		}
	}

	return nil
}
