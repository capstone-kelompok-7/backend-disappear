package faker

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
	"time"
)

func CarouselFaker(db *gorm.DB) []*entities.CarouselModels {
	carousels := []*entities.CarouselModels{
		{
			ID:        1,
			Name:      "Lets Plant Tree",
			Photo:     "https://res.cloudinary.com/dufa4bel6/image/upload/v1702302228/disappear/med4l94zeomooc0ipxlb.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		{
			ID:        2,
			Name:      "Cristmass Day",
			Photo:     "https://res.cloudinary.com/dufa4bel6/image/upload/v1702302550/disappear/pvkmms1w5dvw3o85vvqf.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		{
			ID:        3,
			Name:      "Best Product",
			Photo:     "https://res.cloudinary.com/dufa4bel6/image/upload/v1702302502/disappear/j4wrujgd20wzqkuhdsig.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	return carousels
}
