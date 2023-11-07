package database

import (
	products "github.com/capstone-kelompok-7/backend-disappear/module/product/domain"
	users "github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(users.UserModels{}, users.AddressModels{}, products.Category{}, products.Product{}, products.ProductPhotos{}, products.Review{})
	if err != nil {
		return
	}
}
