package service

import (
	"errors"
	"testing"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant"
	assistantMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/mocks"
	assistants "github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/service"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	productsMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/product/mocks"
	products "github.com/capstone-kelompok-7/backend-disappear/module/feature/product/service"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestService(t *testing.T) (*mocks.RepositoryReviewInterface, review.ServiceReviewInterface, product.ServiceProductInterface, assistant.ServiceAssistantInterface) {
	repo := mocks.NewRepositoryReviewInterface(t)
	repoProduct := productsMocks.NewRepositoryProductInterface(t)
	repoAssistant := assistantMocks.NewRepositoryAssistantInterface(t)

	assistantService := assistants.NewAssistantService(repoAssistant, nil, config.Config{})
	productService := products.NewProductService(repoProduct, assistantService)
	reviewService := NewReviewService(repo, productService)

	return repo, reviewService, productService, assistantService
}

func TestCreateReview(t *testing.T) {
	repo, reviewService, _, _ := setupTestService(t)

	reviewData := &entities.ReviewModels{
		UserID:      2,
		ProductID:   2,
		Rating:      4,
		Description: "Great product!",
		Date:        time.Now(),
		CreatedAt:   time.Now(),
	}

	productServiceMock := productsMocks.NewServiceProductInterface(t)
	product := &entities.ProductModels{
		ID:    2,
		Name:  "Product ABC",
		Price: 1000,
	}

	createdReview := &entities.ReviewModels{
		ID:          1,
		UserID:      reviewData.UserID,
		ProductID:   reviewData.ProductID,
		Rating:      reviewData.Rating,
		Description: reviewData.Description,
		Date:        reviewData.Date,
		CreatedAt:   reviewData.CreatedAt,
	}

	t.Run("Success Case - Create Review", func(t *testing.T) {
		productServiceMock.On("GetProductByID", reviewData.ProductID).Return(product, nil)
		reviewService.(*ReviewService).productService = productServiceMock

		repo.On("CreateReview", mock.AnythingOfType("*entities.ReviewModels")).Return(createdReview, nil)
		productServiceMock.On("UpdateTotalReview", reviewData.ProductID).Return(nil)
		repo.On("CountAverageRating", reviewData.ProductID).Return(4.5, nil)
		productServiceMock.On("UpdateProductRating", reviewData.ProductID, 4.5).Return(nil)
		result, err := reviewService.CreateReview(reviewData)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdReview, result)
		productServiceMock.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		productServiceMock := productsMocks.NewServiceProductInterface(t)
		productServiceMock.On("GetProductByID", reviewData.ProductID).Return(nil, errors.New("Product not found"))
		reviewService.(*ReviewService).productService = productServiceMock

		_, err := reviewService.CreateReview(reviewData)

		assert.Error(t, err)
		assert.EqualError(t, err, "produk tidak ditemukan")
		productServiceMock.AssertExpectations(t)
	})

	t.Run("InvalidRating", func(t *testing.T) {
		reviewData := &entities.ReviewModels{
			UserID:      2,
			ProductID:   2,
			Rating:      6,
			Description: "Great product!",
			Date:        time.Now(),
			CreatedAt:   time.Now(),
		}

		productServiceMock := productsMocks.NewServiceProductInterface(t)
		productServiceMock.On("GetProductByID", reviewData.ProductID).Return(&entities.ProductModels{}, nil)
		reviewService.(*ReviewService).productService = productServiceMock

		_, err := reviewService.CreateReview(reviewData)

		assert.Error(t, err)
		assert.EqualError(t, err, "rating tidak boleh lebih dari 5")
		productServiceMock.AssertExpectations(t)
	})

	t.Run("Failed Case - Update Total Reviews", func(t *testing.T) {
		reviewData := &entities.ReviewModels{
			UserID:      2,
			ProductID:   2,
			Rating:      4,
			Description: "Great product!",
			Date:        time.Now(),
			CreatedAt:   time.Now(),
		}

		productServiceMock := productsMocks.NewServiceProductInterface(t)
		product := &entities.ProductModels{
			ID:    2,
			Name:  "Product ABC",
			Price: 1000,
		}
		productServiceMock.On("GetProductByID", reviewData.ProductID).Return(product, nil)
		reviewService.(*ReviewService).productService = productServiceMock

		createdReview := &entities.ReviewModels{
			ID:          1,
			UserID:      reviewData.UserID,
			ProductID:   reviewData.ProductID,
			Rating:      reviewData.Rating,
			Description: reviewData.Description,
			Date:        reviewData.Date,
			CreatedAt:   reviewData.CreatedAt,
		}
		repo.On("CreateReview", mock.AnythingOfType("*entities.ReviewModels")).Return(createdReview, nil)
		productServiceMock.On("UpdateTotalReview", reviewData.ProductID).Return(errors.New("gagal memperbarui total reviews"))
		_, err := reviewService.CreateReview(reviewData)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal memperbarui total reviews")

		productServiceMock.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Count Average Rating", func(t *testing.T) {
		repo, reviewService, _, _ := setupTestService(t)
		productServiceMock.On("GetProductByID", reviewData.ProductID).Return(product, nil)
		reviewService.(*ReviewService).productService = productServiceMock
		repo.On("CreateReview", mock.AnythingOfType("*entities.ReviewModels")).Return(createdReview, nil)
		productServiceMock.On("UpdateTotalReview", reviewData.ProductID).Return(nil)
		repo.On("CountAverageRating", reviewData.ProductID).Return(0.0, errors.New("gagal menghitung rata-rata rating produk"))

		_, err := reviewService.CreateReview(reviewData)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal menghitung rata-rata rating produk")

		productServiceMock.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Update Product Rating", func(t *testing.T) {
		reviewData := &entities.ReviewModels{
			UserID:      2,
			ProductID:   2,
			Rating:      4,
			Description: "Great product!",
			Date:        time.Now(),
			CreatedAt:   time.Now(),
		}

		productServiceMock := productsMocks.NewServiceProductInterface(t)
		product := &entities.ProductModels{
			ID:    2,
			Name:  "Product ABC",
			Price: 1000,
		}
		productServiceMock.On("GetProductByID", reviewData.ProductID).Return(product, nil)
		reviewService.(*ReviewService).productService = productServiceMock

		createdReview := &entities.ReviewModels{
			ID:          1,
			UserID:      reviewData.UserID,
			ProductID:   reviewData.ProductID,
			Rating:      reviewData.Rating,
			Description: reviewData.Description,
			Date:        reviewData.Date,
			CreatedAt:   reviewData.CreatedAt,
		}

		repo.On("CreateReview", mock.AnythingOfType("*entities.ReviewModels")).Return(createdReview, nil)
		productServiceMock.On("UpdateTotalReview", reviewData.ProductID).Return(nil)
		repo.On("CountAverageRating", reviewData.ProductID).Return(4.5, nil)
		productServiceMock.On("UpdateProductRating", reviewData.ProductID, 4.5).Return(errors.New("gagal memperbarui rating produk"))

		_, err := reviewService.CreateReview(reviewData)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal memperbarui rating produk")

		productServiceMock.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Create Review", func(t *testing.T) {
		repo, reviewService, _, _ := setupTestService(t)
		productServiceMock := productsMocks.NewServiceProductInterface(t)
		productServiceMock.On("GetProductByID", reviewData.ProductID).Return(&entities.ProductModels{}, nil)
		reviewService.(*ReviewService).productService = productServiceMock

		expectedErr := errors.New("failed to create review")
		repo.On("CreateReview", mock.AnythingOfType("*entities.ReviewModels")).Return(nil, expectedErr)

		_, err := reviewService.CreateReview(reviewData)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		productServiceMock.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestCreateReviewImages(t *testing.T) {
	repo, reviewService, _, _ := setupTestService(t)

	reviewData := &entities.ReviewPhotoModels{
		ReviewID:  1,
		ImageURL:  "https://example.com/image.jpg",
		CreatedAt: time.Now(),
	}

	createdReviewPhoto := &entities.ReviewPhotoModels{
		ReviewID:  reviewData.ReviewID,
		ImageURL:  reviewData.ImageURL,
		CreatedAt: reviewData.CreatedAt,
	}

	t.Run("Success Case - Create Review Images", func(t *testing.T) {
		repo.On("GetReviewsById", reviewData.ReviewID).Return(&entities.ReviewModels{}, nil)
		repo.On("CreateReviewImages", mock.AnythingOfType("*entities.ReviewPhotoModels")).Return(createdReviewPhoto, nil)

		result, err := reviewService.CreateReviewImages(reviewData)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdReviewPhoto, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Produk Not Found", func(t *testing.T) {
		repo, reviewService, _, _ := setupTestService(t)
		repo.On("GetReviewsById", reviewData.ReviewID).Return(nil, errors.New("produk tidak ditemukan"))

		_, err := reviewService.CreateReviewImages(reviewData)

		assert.Error(t, err)
		assert.EqualError(t, err, "produk tidak ditemukan")
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Create Review Images", func(t *testing.T) {
		repo, reviewService, _, _ := setupTestService(t)
		repo.On("GetReviewsById", reviewData.ReviewID).Return(&entities.ReviewModels{}, nil)

		expectedErr := errors.New("failed to create review images")
		repo.On("CreateReviewImages", mock.AnythingOfType("*entities.ReviewPhotoModels")).Return(nil, expectedErr)

		_, err := reviewService.CreateReviewImages(reviewData)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})
}

func TestGetReviewById(t *testing.T) {
	repo, reviewService, _, _ := setupTestService(t)

	reviewID := uint64(1)
	review := &entities.ReviewModels{
		ID:          reviewID,
		UserID:      2,
		ProductID:   2,
		Rating:      4,
		Description: "Great product!",
		Date:        time.Now(),
		CreatedAt:   time.Now(),
	}

	t.Run("Success Case - Get Review by ID", func(t *testing.T) {
		repo.On("GetReviewsById", reviewID).Return(review, nil)

		result, err := reviewService.GetReviewById(reviewID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, review, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Review Not Found", func(t *testing.T) {
		repo, reviewService, _, _ := setupTestService(t)
		repo.On("GetReviewsById", reviewID).Return(nil, errors.New("reviews tidak di temukan"))

		_, err := reviewService.GetReviewById(reviewID)

		assert.Error(t, err)
		assert.EqualError(t, err, "reviews tidak di temukan")
		repo.AssertExpectations(t)
	})
}

func TestCountAverageRating(t *testing.T) {
	repo, reviewService, _, _ := setupTestService(t)

	productID := uint64(1)
	averageRating := 4.5

	t.Run("Success Case - Count Average Rating", func(t *testing.T) {
		repo.On("CountAverageRating", productID).Return(averageRating, nil)

		result, err := reviewService.CountAverageRating(productID)

		assert.Nil(t, err)
		assert.Equal(t, averageRating, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Count Average Rating", func(t *testing.T) {
		repo, reviewService, _, _ := setupTestService(t)
		repo.On("CountAverageRating", productID).Return(0.0, errors.New("gagal menghitung rata - rata rating"))

		_, err := reviewService.CountAverageRating(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal menghitung rata - rata rating")
		repo.AssertExpectations(t)
	})
}

func TestGetReviewsProductByID(t *testing.T) {
	repo, reviewService, _, _ := setupTestService(t)

	productID := uint64(1)
	product := &entities.ProductModels{
		ID:    productID,
		Name:  "Product ABC",
		Price: 1000,
	}

	t.Run("Success Case - Get Reviews by Product ID", func(t *testing.T) {
		repo.On("GetReviewsProductByID", productID).Return(product, nil)

		result, err := reviewService.GetReviewsProductByID(productID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, product, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Product Not Found", func(t *testing.T) {
		repo, reviewService, _, _ := setupTestService(t)
		repo.On("GetReviewsProductByID", productID).Return(nil, errors.New("produk tidak ditemukan"))

		_, err := reviewService.GetReviewsProductByID(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "produk tidak ditemukan")
		repo.AssertExpectations(t)
	})
}
