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

func (r *ReviewRepository) CountAverageRating(productID uint64) (float64, error) {
	var averageRating float64

	query := "SELECT AVG(rating) FROM reviews WHERE product_id = ?"
	if err := r.db.Raw(query, productID).Scan(&averageRating).Error; err != nil {
		return 0, err
	}

	return averageRating, nil
}

func (r *ReviewRepository) GetReviewsProductByID(productID uint64) (*entities.ProductModels, error) {
	var product *entities.ProductModels

	if err := r.db.
		Preload("ProductReview").Preload("ProductReview.User").Preload("ProductReview.Photos").
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}
