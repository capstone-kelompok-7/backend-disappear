package database

import (
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(config config.Config) *gorm.DB {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("Database : cannot connect to database ", err.Error())
		return nil
	}
	logrus.Info("Database connected successfully")
	return db
}
