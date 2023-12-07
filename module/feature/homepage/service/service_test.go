package service

import (
	"errors"
	"testing"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/homepage/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetBestSellingProducts(t *testing.T) {
	repo := mocks.NewRepositoryHomepageInterface(t)
	service := NewHomepageService(repo)

	expectedResult := []*entities.ProductModels{
		{
			ID:          1,
			Name:        "Best",
			Description: "Best selling product",
		},
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetBestSellingProducts", 5).Return(expectedResult, nil)

		result, err := service.GetBestSellingProducts(5)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - ErrorFromRepository", func(t *testing.T) {
		repo := mocks.NewRepositoryHomepageInterface(t)
		service := NewHomepageService(repo)
		expectedError := errors.New("repository error")
		repo.On("GetBestSellingProducts", 5).Return(nil, expectedError)

		result, err := service.GetBestSellingProducts(5)

		assert.Error(t, err, "expected an error from the repository")
		assert.Nil(t, result, "expected result to be nil")
		repo.AssertExpectations(t)
	})
}

func TestGetCategory(t *testing.T) {
	repo := mocks.NewRepositoryHomepageInterface(t)
	service := NewHomepageService(repo)

	expectedResult := []*entities.CategoryModels{
		{
			ID:   1,
			Name: "Category 1",
		},
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetFiveCategories").Return(expectedResult, nil)

		result, err := service.GetCategory()

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - ErrorFromRepository", func(t *testing.T) {
		repo := mocks.NewRepositoryHomepageInterface(t)
		service := NewHomepageService(repo)

		expectedError := errors.New("repository error")
		repo.On("GetFiveCategories").Return(nil, expectedError)

		result, err := service.GetCategory()

		assert.Error(t, err, "expected an error from the repository")
		assert.Nil(t, result, "expected result to be nil")
		repo.AssertExpectations(t)
	})
}

func TestGetCarousel(t *testing.T) {
	repo := mocks.NewRepositoryHomepageInterface(t)
	service := NewHomepageService(repo)

	expectedResult := []*entities.CarouselModels{
		{
			ID:   1,
			Name: "Carousel 1",
		},
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetFiveCarousel").Return(expectedResult, nil)

		result, err := service.GetCarousel()

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - ErrorFromRepository", func(t *testing.T) {
		repo := mocks.NewRepositoryHomepageInterface(t)
		service := NewHomepageService(repo)

		expectedError := errors.New("repository error")
		repo.On("GetFiveCarousel").Return(nil, expectedError)

		result, err := service.GetCarousel()

		assert.Error(t, err, "expected an error from the repository")
		assert.Nil(t, result, "expected result to be nil")
		repo.AssertExpectations(t)
	})
}

func TestGetChallenge(t *testing.T) {
	repo := mocks.NewRepositoryHomepageInterface(t)
	service := NewHomepageService(repo)

	expectedResult := []*entities.ChallengeModels{
		{
			ID:    1,
			Title: "Challenge 1",
		},
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetFiveChallenge").Return(expectedResult, nil)

		result, err := service.GetChallenge()

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - ErrorFromRepository", func(t *testing.T) {
		repo := mocks.NewRepositoryHomepageInterface(t)
		service := NewHomepageService(repo)

		expectedError := errors.New("repository error")
		repo.On("GetFiveChallenge").Return(nil, expectedError)

		result, err := service.GetChallenge()

		assert.Error(t, err, "expected an error from the repository")
		assert.Nil(t, result, "expected result to be nil")
		repo.AssertExpectations(t)
	})
}

func TestGetArticle(t *testing.T) {
	repo := mocks.NewRepositoryHomepageInterface(t)
	service := NewHomepageService(repo)

	expectedResult := []*entities.ArticleModels{
		{
			ID:    1,
			Title: "Article 1",
		},
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetThreeArticle").Return(expectedResult, nil)

		result, err := service.GetArticle()

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - ErrorFromRepository", func(t *testing.T) {
		repo := mocks.NewRepositoryHomepageInterface(t)
		service := NewHomepageService(repo)

		expectedError := errors.New("repository error")
		repo.On("GetThreeArticle").Return(nil, expectedError)

		result, err := service.GetArticle()

		assert.Error(t, err, "expected an error from the repository")
		assert.Nil(t, result, "expected result to be nil")
		repo.AssertExpectations(t)
	})
}
