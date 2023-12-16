package faker

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
	"time"
)

func ChallengeFaker(db *gorm.DB) []*entities.ChallengeModels {
	challenges := []*entities.ChallengeModels{
		{
			ID:          1,
			Title:       "Tantangan Menanam Pohon",
			Photo:       "https://res.cloudinary.com/dufa4bel6/image/upload/v1702565839/disappear/dgpziu8cppajpywuki5u.jpg",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 0, 30),
			Description: "Challenge untuk melakukan daur ulang selama 30 hari berturut-turut...",
			Status:      "Belum Kadaluwarsa",
			Exp:         50,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
		{
			ID:          2,
			Title:       "Challenge Kurangi Limbah Plastik",
			Photo:       "https://res.cloudinary.com/dufa4bel6/image/upload/v1702566045/disappear/bvw6bov7hwkqrib5igaa.jpg",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 0, 45),
			Description: "Challenge untuk mengurangi penggunaan limbah plastik dalam 45 hari...",
			Status:      "Kadaluwarsa",
			Exp:         70,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
	}

	return challenges
}
