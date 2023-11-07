package database

import (
	products "github.com/capstone-kelompok-7/backend-disappear/module/product/domain"
	review "github.com/capstone-kelompok-7/backend-disappear/module/review/domain"
	users "github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	voucher "github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(voucher.VoucherModels{}, users.UserModels{}, users.AddressModels{}, products.Category{}, products.Product{}, products.ProductPhotos{}, review.Review{})
	if err != nil {
		return
	}
}
