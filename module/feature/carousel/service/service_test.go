package service

import (
	"errors"
	"testing"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCarouselService_CalculatePaginationValues(t *testing.T) {
	service := &CarouselService{}

	t.Run("Page less than or equal to zero should default to 1", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(0, 100, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page exceeds total pages should set to total pages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(15, 100, 8)

		assert.Equal(t, 13, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page within limits should return correct values", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(2, 100, 8)

		assert.Equal(t, 2, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Total items not perfectly divisible by perPage should round totalPages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(1, 95, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 12, totalPages)
	})
}

func TestCarouselService_GetNextPage(t *testing.T) {
	service := &CarouselService{}

	t.Run("Next Page Within Total Pages", func(t *testing.T) {
		currentPage := 3
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, currentPage+1, nextPage)
	})

	t.Run("Next Page Equal to Total Pages", func(t *testing.T) {
		currentPage := 5
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, totalPages, nextPage)
	})
}

func TestCarouselService_GetPrevPage(t *testing.T) {
	service := &CarouselService{}

	t.Run("Previous Page Within Bounds", func(t *testing.T) {
		currentPage := 3

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage-1, prevPage)
	})

	t.Run("Previous Page at Lower Bound", func(t *testing.T) {
		currentPage := 1

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage, prevPage)
	})
}

func TestCarouselService_GetAll(t *testing.T) {
	repo := mocks.NewRepositoryCarouselInterface(t)
	service := NewCarouselService(repo)

	carousels := []*entities.CarouselModels{
		{ID: 1, Name: "Carousel 1", Photo: "carousel1.jpg"},
		{ID: 2, Name: "Carousel 2", Photo: "carousel2.jpg"},
	}

	t.Run("Success Case - Carousel Found", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("FindAll", 1, 10).Return(carousels, nil).Once()
		repo.On("GetTotalCarouselCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetAll(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(carousels), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetTotalCarouselCount Error", func(t *testing.T) {
		expectedErr := errors.New(" GetTotalCarouselCount Error")
		repo.On("FindAll", 1, 10).Return(carousels, nil).Once()
		repo.On("GetTotalCarouselCount").Return(int64(0), expectedErr).Once()

		carousels, totalItems, err := service.GetAll(1, 10)

		assert.Error(t, err)
		assert.Nil(t, carousels)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetAll Error", func(t *testing.T) {
		expectedErr := errors.New("FindAll Error")
		repo.On("FindAll", 1, 10).Return(nil, expectedErr).Once()

		carousels, totalItems, err := service.GetAll(1, 10)

		assert.Error(t, err)
		assert.Nil(t, carousels)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})
}

func TestCarouselService_GetCarouselsByName(t *testing.T) {
	repo := mocks.NewRepositoryCarouselInterface(t)
	service := NewCarouselService(repo)

	carousels := []*entities.CarouselModels{
		{ID: 1, Name: "Carousel 1", Photo: "carousel1.jpg"},
		{ID: 2, Name: "Carousel 2", Photo: "carousel2.jpg"},
	}
	name := "Test"

	t.Run("Success Case - Carousels Found by Name", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("FindByName", 1, 10, name).Return(carousels, nil).Once()
		repo.On("GetTotalCarouselCountByName", name).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetCarouselsByName(1, 10, name)

		assert.NoError(t, err)
		assert.Equal(t, len(carousels), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Carousels by Name", func(t *testing.T) {
		expectedErr := errors.New("failed to find carousels by name")
		repo.On("FindByName", 1, 10, name).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetCarouselsByName(1, 10, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Carousel Count by Name", func(t *testing.T) {
		expectedErr := errors.New("failed to get total carousel count by name")
		repo.On("FindByName", 1, 10, name).Return(carousels, nil).Once()
		repo.On("GetTotalCarouselCountByName", name).Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetCarouselsByName(1, 10, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestCarouselService_GetCarouselById(t *testing.T) {
	repo := mocks.NewRepositoryCarouselInterface(t)
	service := NewCarouselService(repo)

	carousels := &entities.CarouselModels{
		ID:    1,
		Name:  "Carousel 1",
		Photo: "carousel2.jpg",
	}

	expectedCarousels := &entities.CarouselModels{
		ID:    carousels.ID,
		Name:  carousels.Name,
		Photo: carousels.Photo,
	}

	t.Run("Success Case - Carousel Found", func(t *testing.T) {
		carouselID := uint64(1)
		repo.On("GetCarouselById", carouselID).Return(expectedCarousels, nil).Once()

		result, err := service.GetCarouselById(carouselID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCarousels.ID, result.ID)
		assert.Equal(t, expectedCarousels.Name, result.Name)
		assert.Equal(t, expectedCarousels.Photo, result.Photo)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Carousel Not Found", func(t *testing.T) {
		carouselID := uint64(2)

		expectedErr := errors.New("carousel tidak ditemukan")
		repo.On("GetCarouselById", carouselID).Return(nil, expectedErr).Once()

		result, err := service.GetCarouselById(carouselID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})

}

func TestCarouselService_CreateCarousel(t *testing.T) {
	repo := mocks.NewRepositoryCarouselInterface(t)
	service := NewCarouselService(repo)

	carousels := &entities.CarouselModels{
		Name:  "Carousel 1",
		Photo: "carousel2.jpg",
	}

	expectedCarousels := &entities.CarouselModels{
		Name:  carousels.Name,
		Photo: carousels.Photo,
	}

	t.Run("Success Case", func(t *testing.T) {

		repo.On("CreateCarousel", carousels).Return(expectedCarousels, nil).Once()

		result, err := service.CreateCarousel(carousels)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCarousels.ID, result.ID)
		assert.Equal(t, expectedCarousels.Name, result.Name)
		assert.Equal(t, expectedCarousels.Photo, result.Photo)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case", func(t *testing.T) {

		expectedErr := errors.New("CreateCarousel")
		repo.On("CreateCarousel", carousels).Return(nil, expectedErr).Once()

		result, err := service.CreateCarousel(carousels)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})

}

func TestCarouselService_UpdateCarousel(t *testing.T) {
	repo := mocks.NewRepositoryCarouselInterface(t)
	service := NewCarouselService(repo)

	carousels := &entities.CarouselModels{
		ID:    1,
		Name:  "Carousel 1",
		Photo: "carousel2.jpg",
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetCarouselById", carousels.ID).Return(carousels, nil).Once()
		repo.On("UpdateCarousel", carousels.ID, carousels).Return(nil).Once()

		result := service.UpdateCarousel(carousels.ID, carousels)

		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Carousel Not Found", func(t *testing.T) {
		expectedErr := errors.New("carousel tidak ditemukan")
		repo.On("GetCarouselById", carousels.ID).Return(nil, expectedErr).Once()

		result := service.UpdateCarousel(carousels.ID, carousels)

		assert.Error(t, result)
		assert.Equal(t, expectedErr, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Carousel Found But Update Failed", func(t *testing.T) {
		expectedUpdateErr := errors.New("UpdateCarousel")
		repo.On("GetCarouselById", carousels.ID).Return(carousels, nil).Once()
		repo.On("UpdateCarousel", carousels.ID, carousels).Return(expectedUpdateErr).Once()

		result := service.UpdateCarousel(carousels.ID, carousels)

		assert.Error(t, result)
		assert.Equal(t, expectedUpdateErr, result)
		repo.AssertExpectations(t)
	})

}

func TestCarouselService_DeleteCarousel(t *testing.T) {
	repo := mocks.NewRepositoryCarouselInterface(t)
	service := NewCarouselService(repo)
	carousels := &entities.CarouselModels{
		ID:    1,
		Name:  "Carousel 1",
		Photo: "carousel2.jpg",
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetCarouselById", carousels.ID).Return(carousels, nil).Once()
		repo.On("DeleteCarousel", carousels.ID).Return(nil).Once()

		err := service.DeleteCarousel(carousels.ID)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Carousel Not Found", func(t *testing.T) {
		expectedErr := errors.New("carousel tidak ditemukan")
		repo.On("GetCarouselById", carousels.ID).Return(nil, expectedErr).Once()

		result := service.DeleteCarousel(carousels.ID)

		assert.Error(t, result)
		assert.Equal(t, expectedErr, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Carousel Found But Delete Failed", func(t *testing.T) {
		expectedDeleteErr := errors.New("DeleteCarousel")
		repo.On("GetCarouselById", carousels.ID).Return(carousels, nil).Once()
		repo.On("DeleteCarousel", carousels.ID).Return(expectedDeleteErr).Once()

		result := service.DeleteCarousel(carousels.ID)

		assert.Error(t, result)
		assert.Equal(t, expectedDeleteErr, result)
		repo.AssertExpectations(t)
	})

}
