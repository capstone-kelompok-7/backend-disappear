package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review"
	"time"
)

type ReviewService struct {
	repo review.RepositoryReviewInterface
}

func NewReviewService(reviewRepo review.RepositoryReviewInterface) review.ServiceReviewInterface {
	return &ReviewService{
		repo: reviewRepo,
	}
}

func (s *ReviewService) CreateReview(reviewData *entities.ReviewModels) (*entities.ReviewModels, error) {

	if reviewData.Rating > 5 {
		return nil, errors.New("rating tidak boleh lebih dari 5")
	}
	value := &entities.ReviewModels{
		UserID:      reviewData.UserID,
		ProductID:   reviewData.ProductID,
		Rating:      reviewData.Rating,
		Description: reviewData.Description,
		CreatedAt:   time.Now(),
	}

	createdReview, err := s.repo.CreateReview(value)
	if err != nil {
		return nil, err
	}

	return createdReview, nil
}
