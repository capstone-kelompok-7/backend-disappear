package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) review.RepositoryReviewInterface {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) CreateReview(newData *entities.ReviewModels) (*entities.ReviewModels, error) {
	err := r.db.Create(&newData).Error
	if err != nil {
		return nil, err
	}
	return newData, nil
}

func (r *ReviewRepository) CreateReviewImages(newData *entities.ReviewPhotoModels) (*entities.ReviewPhotoModels, error) {
	err := r.db.Create(&newData).Error
	if err != nil {
		return nil, err
	}
	return newData, nil
}

func (r *ReviewRepository) GetReviewsById(reviewID uint64) (*entities.ReviewModels, error) {
	var reviews *entities.ReviewModels

	if err := r.db.Preload("Photos").Where("id = ? AND deleted_at IS NULL", reviewID).First(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}
