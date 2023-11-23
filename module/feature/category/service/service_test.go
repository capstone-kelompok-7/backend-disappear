package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCategoryService_CalculatePaginationValues(t *testing.T) {
	service := &CategoryService{}

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

func TestCategoryService_GetNextPage(t *testing.T) {
	service := &CategoryService{}

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

func TestCategoryService_GetPrevPage(t *testing.T) {
	service := &CategoryService{}

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

func TestCategoryService_GetAll(t *testing.T) {
	repo := mocks.NewRepositoryCategoryInterface(t)
	service := NewCategoryService(repo)

	categories := []*entities.CategoryModels{
		{ID: 1, Name: "Category 1", Photo: "categories1.jpg"},
		{ID: 2, Name: "Category 2", Photo: "categories2.jpg"},
	}

	t.Run("Success Case - Category Found", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("FindAll", 1, 10).Return(categories, nil).Once()
		repo.On("GetTotalCategoryCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetAllCategory(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(categories), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetTotalCategoryCount Error", func(t *testing.T) {
		expectedErr := errors.New(" GetTotalCategoryCount Error")
		repo.On("FindAll", 1, 10).Return(categories, nil).Once()
		repo.On("GetTotalCategoryCount").Return(int64(0), expectedErr).Once()

		category, totalItems, err := service.GetAllCategory(1, 10)

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetAll Error", func(t *testing.T) {
		expectedErr := errors.New("FindAll Error")
		repo.On("FindAll", 1, 10).Return(nil, expectedErr).Once()

		category, totalItems, err := service.GetAllCategory(1, 10)

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})
}

func TestCategoryService_GetCategoryByName(t *testing.T) {
	repo := mocks.NewRepositoryCategoryInterface(t)
	service := NewCategoryService(repo)

	categories := []*entities.CategoryModels{
		{ID: 1, Name: "Category 1", Photo: "categories1.jpg"},
		{ID: 2, Name: "Category 2", Photo: "categories2.jpg"},
	}
	name := "Test"

	t.Run("Success Case - Category Found by Name", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("FindByName", 1, 10, name).Return(categories, nil).Once()
		repo.On("GetTotalCategoryCountByName", name).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetCategoryByName(1, 10, name)

		assert.NoError(t, err)
		assert.Equal(t, len(categories), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Category by Name", func(t *testing.T) {
		expectedErr := errors.New("failed to find categories by name")
		repo.On("FindByName", 1, 10, name).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetCategoryByName(1, 10, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Category Count by Name", func(t *testing.T) {
		expectedErr := errors.New("failed to get total carousel count by name")
		repo.On("FindByName", 1, 10, name).Return(categories, nil).Once()
		repo.On("GetTotalCategoryCountByName", name).Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetCategoryByName(1, 10, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestCategoryService_CreateCategory(t *testing.T) {
	repo := mocks.NewRepositoryCategoryInterface(t)
	service := NewCategoryService(repo)

	categories := &entities.CategoryModels{
		Name:  "Category 1",
		Photo: "category.jpg",
	}

	expectedCategories := &entities.CategoryModels{
		Name:  categories.Name,
		Photo: categories.Photo,
	}

	t.Run("Success Case", func(t *testing.T) {

		repo.On("CreateCategory", categories).Return(expectedCategories, nil).Once()

		result, err := service.CreateCategory(categories)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCategories.ID, result.ID)
		assert.Equal(t, expectedCategories.Name, result.Name)
		assert.Equal(t, expectedCategories.Photo, result.Photo)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case", func(t *testing.T) {

		expectedErr := errors.New("CreateCategory")
		repo.On("CreateCategory", categories).Return(nil, expectedErr).Once()

		result, err := service.CreateCategory(categories)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})

}

func TestCategoryService_UpdateCategory(t *testing.T) {
	repo := mocks.NewRepositoryCategoryInterface(t)
	service := NewCategoryService(repo)

	categories := &entities.CategoryModels{
		Name:  "Category 1",
		Photo: "category.jpg",
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetCategoryById", categories.ID).Return(categories, nil).Once()
		repo.On("UpdateCategoryById", categories.ID, categories).Return(nil).Once()

		result := service.UpdateCategoryById(categories.ID, categories)

		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Category Not Found", func(t *testing.T) {
		expectedErr := errors.New("kategori tidak ditemukan")
		repo.On("GetCategoryById", categories.ID).Return(nil, expectedErr).Once()

		result := service.UpdateCategoryById(categories.ID, categories)

		assert.Error(t, result)
		assert.Equal(t, expectedErr, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Category Found But Update Failed", func(t *testing.T) {
		expectedUpdateErr := errors.New("UpdateCategoryById")
		repo.On("GetCategoryById", categories.ID).Return(categories, nil).Once()
		repo.On("UpdateCategoryById", categories.ID, categories).Return(expectedUpdateErr).Once()

		result := service.UpdateCategoryById(categories.ID, categories)

		assert.Error(t, result)
		assert.Equal(t, expectedUpdateErr, result)
		repo.AssertExpectations(t)
	})

}

func TestCategoryService_DeleteCategory(t *testing.T) {
	repo := mocks.NewRepositoryCategoryInterface(t)
	service := NewCategoryService(repo)

	categories := &entities.CategoryModels{
		Name:  "Category 1",
		Photo: "category.jpg",
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetCategoryById", categories.ID).Return(categories, nil).Once()
		repo.On("DeleteCategoryById", categories.ID).Return(nil).Once()

		err := service.DeleteCategoryById(categories.ID)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Category Not Found", func(t *testing.T) {
		expectedErr := errors.New("kategori tidak ditemukan")
		repo.On("GetCategoryById", categories.ID).Return(nil, expectedErr).Once()

		result := service.DeleteCategoryById(categories.ID)

		assert.Error(t, result)
		assert.Equal(t, expectedErr, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Category Found But Delete Failed", func(t *testing.T) {
		expectedDeleteErr := errors.New("DeleteCategoryById")
		repo.On("GetCategoryById", categories.ID).Return(categories, nil).Once()
		repo.On("DeleteCategoryById", categories.ID).Return(expectedDeleteErr).Once()

		result := service.DeleteCategoryById(categories.ID)

		assert.Error(t, result)
		assert.Equal(t, expectedDeleteErr, result)
		repo.AssertExpectations(t)
	})

}

func TestCategoryService_GetCategoryById(t *testing.T) {
	repo := mocks.NewRepositoryCategoryInterface(t)
	service := NewCategoryService(repo)

	categories := &entities.CategoryModels{
		ID:    1,
		Name:  "Category 1",
		Photo: "category2.jpg",
	}

	expectedCategory := &entities.CategoryModels{
		ID:    categories.ID,
		Name:  categories.Name,
		Photo: categories.Photo,
	}

	t.Run("Success Case - Category Found", func(t *testing.T) {
		categoryID := uint64(1)
		repo.On("GetCategoryById", categoryID).Return(expectedCategory, nil).Once()

		result, err := service.GetCategoryById(categoryID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCategory.ID, result.ID)
		assert.Equal(t, expectedCategory.Name, result.Name)
		assert.Equal(t, expectedCategory.Photo, result.Photo)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Category Not Found", func(t *testing.T) {
		categoryID := uint64(2)

		expectedErr := errors.New("kategori tidak ditemukan")
		repo.On("GetCategoryById", categoryID).Return(nil, expectedErr).Once()

		result, err := service.GetCategoryById(categoryID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})

}
