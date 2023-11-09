package database

import (
	entities "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		entities.VoucherModels{},
		entities.UserModels{},
		entities.AddressModels{},
		entities.CategoryModels{},
		entities.ProductModels{},
		entities.ProductPhotosModels{},
		entities.ReviewModels{},
		entities.Articles{},
		entities.OTPModels{},
		entities.ChallengeModels{},
	)

	if err != nil {
		return
	}
}
