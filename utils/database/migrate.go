package database

import (
	users "github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(users.UserModels{}, users.AddressModels{})
	if err != nil {
		return
	}
}
