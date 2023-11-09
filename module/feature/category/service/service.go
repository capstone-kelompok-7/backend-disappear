package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category"
	"math"
)

type CategoryService struct {
	categoryRepo category.RepositoryCategoryInterface
}

func NewCategoryService(categoryRepo category.RepositoryCategoryInterface) category.ServiceCategoryInterface {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(categoryData *entities.CategoryModels) (*entities.CategoryModels, error) {
	createdCategory, err := s.categoryRepo.CreateCategory(categoryData)
	if err != nil {
		return nil, err
	}

	return createdCategory, nil
}

func (s *CategoryService) GetAllCategory(page, perPage int) ([]*entities.CategoryModels, int64, error) {
	categories, err := s.categoryRepo.FindAll(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.categoryRepo.GetTotalCategoryCount()
	if err != nil {
		return nil, 0, err
	}

	return categories, totalItems, nil
}

func (s *CategoryService) GetCategoryByName(page int, perPage int, name string) ([]*entities.CategoryModels, int64, error) {
	categories, err := s.categoryRepo.FindByName(page, perPage, name)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.categoryRepo.GetTotalCategoryCountByName(name)
	if err != nil {
		return nil, 0, err
	}

	return categories, totalItems, nil
}

func (s *CategoryService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if page > totalPages {
		page = totalPages
	}

	return page, totalPages
}

func (s *CategoryService) GetNextPage(currentPage int, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}

	return totalPages
}

func (s *CategoryService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}

	return 1
}

func (s *CategoryService) UpdateCategoryById(id int, updatedCategory *entities.CategoryModels) (*entities.CategoryModels, error) {
	existingCategory, err := s.categoryRepo.GetCategoryById(id)
	if err != nil {
		return nil, err
	}

	if existingCategory == nil {
		return nil, errors.New("Kategori tidak ditemukan")
	}

	updatedCategory, err = s.categoryRepo.UpdateCategoryById(id, updatedCategory)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (s *CategoryService) DeleteCategoryById(id int) error {
	existingCategory, err := s.categoryRepo.GetCategoryById(id)
	if err != nil {
		return err
	}

	if existingCategory == nil {
		return errors.New("Kategori tidak ditemukan")
	}

	err = s.categoryRepo.DeleteCategoryById(id)
	if err != nil {
		return err
	}

	return nil
}
