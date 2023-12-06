package database

import (
	entities "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		entities.VoucherModels{},
		entities.UserModels{},
		entities.ArticleBookmarkModels{},
		entities.AddressModels{},
		entities.CategoryModels{},
		entities.ProductModels{},
		entities.ProductPhotosModels{},
		entities.ReviewModels{},
		entities.ArticleModels{},
		entities.OTPModels{},
		entities.ChallengeModels{},
		entities.CarouselModels{},
		entities.ReviewPhotoModels{},
		entities.CartModels{},
		entities.CartItemModels{},
		entities.ChallengeFormModels{},
		entities.OrderModels{},
		entities.OrderDetailsModels{},
		entities.VoucherClaimModels{},
		entities.EnvironmentIssuesModels{},
	)

	if err != nil {
		return
	}
}
