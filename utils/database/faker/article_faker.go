package faker

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
	"time"
)

func ArticleFaker(db *gorm.DB) []*entities.ArticleModels {
	articles := []*entities.ArticleModels{
		{
			ID:        1,
			Title:     "Berapa Banyak Sampah Plastik yang Ada di Lautan?",
			Photo:     "https://res.cloudinary.com/dufa4bel6/image/upload/v1702447426/disappear/zkchk4spud0syegzmpcn.png",
			Content:   "Peneliti memperkirakan jumlah sampah plastik yang memasuki lingkungan perairan dapat meningkat 2,6 kali lipat dari 2016 hingga 2040. Jika tren tersebut terus berlanjut, itu akan berpotensi menjadi hal yang lebih buruk...",
			Author:    "Admin Disappear",
			Views:     0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		{
			ID:        2,
			Title:     "80 Persen Sampah di Laut adalah Sampah dari Daratan",
			Photo:     "https://res.cloudinary.com/dufa4bel6/image/upload/v1702566045/disappear/bvw6bov7hwkqrib5igaa.jpg",
			Content:   "Kepala Badan Riset dan Sumber Daya Manusia Kelautan dan Perikanan (BRSDM) I Nyoman Radiarta mengatakan, sampah laut dan dampak pencemaran terhadap laut sudah menjadi isu skala lokal, nasional, hingga global. Sampah laut atau marine debris dinilai sangat berdampak buruk bagi lingkungan dan biota laut...",
			Author:    "Admin Disappear",
			Views:     0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		{
			ID:        3,
			Title:     "Ayo bepartisipasi dalam menciptakan lingkungan yang aman dan nyaman.",
			Photo:     "https://res.cloudinary.com/dufa4bel6/image/upload/v1702567852/disappear/rp4urjpfzcug1sav2cuz.jpg",
			Content:   "Kepala Badan Riset dan Sumber Daya Manusia Kelautan dan Perikanan (BRSDM) I Nyoman Radiarta mengatakan, sampah laut dan dampak pencemaran terhadap laut sudah menjadi isu skala lokal, nasional, hingga global. Sampah laut atau marine debris dinilai sangat berdampak buruk bagi lingkungan dan biota laut...",
			Author:    "Admin Disappear",
			Views:     0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	return articles
}
