package database

import (
	article "github.com/capstone-kelompok-7/backend-disappear/module/article/domain"
	products "github.com/capstone-kelompok-7/backend-disappear/module/product/domain"
	users "github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	voucher "github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(voucher.VoucherModels{}, users.UserModels{}, users.AddressModels{}, products.Category{}, products.Product{}, products.ProductPhotos{}, products.Review{}, article.Articles{})
	if err != nil {
		return
	}
}
