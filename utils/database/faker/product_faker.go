package faker

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
	"time"
)

func ProductFaker(db *gorm.DB) []*entities.ProductModels {
	products := []*entities.ProductModels{
		{
			ID:          1,
			Name:        "Paperbag",
			Description: "Paperbag yang serbaguna dan ramah lingkungan, cocok untuk berbagai keperluan.",
			GramPlastic: 20,
			Price:       35000,
			Stock:       50,
			Discount:    5000,
			Exp:         30,
			Rating:      5,
			TotalReview: 2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
		{
			ID:          2,
			Name:        "Totebag",
			Description: "Totebag yang bergaya dan tahan lama, ideal untuk penggunaan sehari-hari.",
			GramPlastic: 150,
			Price:       75000,
			Stock:       30,
			Discount:    1500,
			Exp:         60,
			Rating:      0,
			TotalReview: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
		{
			ID:          3,
			Name:        "Stainless Straw",
			Description: "Stainless straw yang dapat digunakan ulang dan tahan karat, sempurna untuk mengurangi limbah plastik.",
			GramPlastic: 120,
			Price:       60000,
			Stock:       40,
			Discount:    1200,
			Exp:         45,
			Rating:      0,
			TotalReview: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
		{
			ID:          4,
			Name:        "Botol Minum",
			Description: "Botol minum yang tahan lama dan anti bocor untuk kebutuhan hidrasi di mana saja.",
			GramPlastic: 200,
			Price:       90000,
			Stock:       20,
			Discount:    20000,
			Exp:         90,
			Rating:      0,
			TotalReview: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
		{
			ID:          5,
			Name:        "Tote Bag",
			Description: "Tote bag yang luas dan terbuat dari bahan berkualitas tinggi, sempurna untuk penggunaan sehari-hari.",
			GramPlastic: 180,
			Price:       80000,
			Stock:       25,
			Discount:    1800,
			Exp:         50,
			Rating:      0,
			TotalReview: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
	}
	return products
}
