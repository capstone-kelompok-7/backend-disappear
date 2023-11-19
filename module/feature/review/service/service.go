package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review"
	"time"
)

type ReviewService struct {
	repo           review.RepositoryReviewInterface
	productService product.ServiceProductInterface
}

func NewReviewService(reviewRepo review.RepositoryReviewInterface, productService product.ServiceProductInterface) review.ServiceReviewInterface {
	return &ReviewService{
		repo:           reviewRepo,
		productService: productService,
	}
}

func (s *ReviewService) CreateReview(reviewData *entities.ReviewModels) (*entities.ReviewModels, error) {
	_, err := s.productService.GetProductByID(reviewData.ProductID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
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

	err = s.productService.UpdateTotalReview(reviewData.ProductID)
	if err != nil {
		return nil, errors.New("gagal memperbarui total reviews")
	}

	return createdReview, nil
}

func (s *ReviewService) CreateReviewImages(reviewData *entities.ReviewPhotoModels) (*entities.ReviewPhotoModels, error) {
	_, err := s.repo.GetReviewsById(reviewData.ReviewID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	value := &entities.ReviewPhotoModels{
		ReviewID:  reviewData.ReviewID,
		ImageURL:  reviewData.ImageURL,
		CreatedAt: time.Now(),
	}

	createdReviewPhoto, err := s.repo.CreateReviewImages(value)
	if err != nil {
		return nil, err
	}

	return createdReviewPhoto, nil
}

func (s *ReviewService) GetReviewById(reviewID uint64) (*entities.ReviewModels, error) {
	result, err := s.repo.GetReviewsById(reviewID)
	if err != nil {
		return nil, errors.New("reviews tidak di temukan")
	}
	return result, nil
}
