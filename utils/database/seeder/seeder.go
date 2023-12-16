package seeder

import (
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

func DBSeed(db *gorm.DB) error {
	for _, seeder := range RegisterSeeders(db) {
		err := db.Debug().Create(seeder.Seeder).Error
		if err != nil {
			return err
		}
	}

	return nil
}
