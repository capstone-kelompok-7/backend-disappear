package review

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryReviewInterface interface {
	CreateReview(newData *entities.ReviewModels) (*entities.ReviewModels, error)
	CreateReviewImages(newData *entities.ReviewPhotoModels) (*entities.ReviewPhotoModels, error)
	GetReviewsById(reviewID uint64) (*entities.ReviewModels, error)
	CountAverageRating(productID uint64) (float64, error)
	GetReviewsProductByID(productID uint64) (*entities.ProductModels, error)
}

type ServiceReviewInterface interface {
	CreateReview(newData *entities.ReviewModels) (*entities.ReviewModels, error)
	CreateReviewImages(reviewData *entities.ReviewPhotoModels) (*entities.ReviewPhotoModels, error)
	GetReviewById(reviewID uint64) (*entities.ReviewModels, error)
	CountAverageRating(productID uint64) (float64, error)
	GetReviewsProductByID(productID uint64) (*entities.ProductModels, error)
}

type HandlerReviewInterface interface {
	CreateReview() echo.HandlerFunc
	CreateReviewImages() echo.HandlerFunc
	GetReviewById() echo.HandlerFunc
	GetDetailReviewProduct() echo.HandlerFunc
}
