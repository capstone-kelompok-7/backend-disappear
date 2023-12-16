package faker

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
	"time"
)

func CategoryFaker(db *gorm.DB) []*entities.CategoryModels {
	categories := []*entities.CategoryModels{
		{
			ID:           1,
			Name:         "Tas",
			Photo:        "https://res.cloudinary.com/dufa4bel6/image/upload/v1702742274/disappear/Tas_gi3ogw.png",
			TotalProduct: 1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    nil,
		},
		{
			ID:           2,
			Name:         "Alat Rumah Tangga",
			Photo:        "https://res.cloudinary.com/dufa4bel6/image/upload/v1702742273/disappear/alat_makan_garpu_mka1dl.png",
			TotalProduct: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    nil,
		},
		{
			ID:           3,
			Name:         "Alat Makan",
			Photo:        "https://res.cloudinary.com/dufa4bel6/image/upload/v1702742274/disappear/Alat_Makan_lvbvb9.png",
			TotalProduct: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    nil,
		},
		{
			ID:           4,
			Name:         "Alat Minum",
			Photo:        "https://res.cloudinary.com/dufa4bel6/image/upload/v1702742274/disappear/Alat_Minum_vhajpp.png",
			TotalProduct: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    nil,
		},
		{
			ID:           5,
			Name:         "Botol Minum",
			Photo:        "https://res.cloudinary.com/dufa4bel6/image/upload/v1702742273/disappear/botol_minum_gwih8j.png",
			TotalProduct: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    nil,
		},
	}
	return categories
}
