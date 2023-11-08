package database

import (
	article "github.com/capstone-kelompok-7/backend-disappear/module/article/domain"
	products "github.com/capstone-kelompok-7/backend-disappear/module/product/domain"
	review "github.com/capstone-kelompok-7/backend-disappear/module/review/domain"
	users "github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	voucher "github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(voucher.VoucherModels{}, users.UserModels{}, users.AddressModels{}, products.CategoryModels{}, products.ProductModels{}, products.ProductPhotosModels{}, review.ReviewModels{}, article.Articles{})
	if err != nil {
		return
	}
}
