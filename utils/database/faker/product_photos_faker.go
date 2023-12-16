package faker

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
	"time"
)

func ProductPhotoFaker(db *gorm.DB) []*entities.ProductPhotosModels {
	productPhotos := []*entities.ProductPhotosModels{
		{
			ID:        1,
			ProductID: 1,
			ImageURL:  "https://res.cloudinary.com/dufa4bel6/image/upload/v1702741226/disappear/fdmomgqfqlpbaeknxccg.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		{
			ID:        2,
			ProductID: 2,
			ImageURL:  "https://res.cloudinary.com/dufa4bel6/image/upload/v1702740398/disappear/j55cubkfhqladvfujkrg.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		{
			ID:        3,
			ProductID: 3,
			ImageURL:  "https://res.cloudinary.com/dufa4bel6/image/upload/v1702741226/disappear/peyrh1csm3nfi5yj2obi.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		{
			ID:        4,
			ProductID: 4,
			ImageURL:  "https://res.cloudinary.com/dufa4bel6/image/upload/v1702741226/disappear/peyrh1csm3nfi5yj2obi.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		{
			ID:        5,
			ProductID: 5,
			ImageURL:  "https://res.cloudinary.com/dufa4bel6/image/upload/v1702741226/disappear/fdmomgqfqlpbaeknxccg.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	return productPhotos
}
