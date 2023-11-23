package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category"
	"math"
)

type CategoryService struct {
	repo category.RepositoryCategoryInterface
}

func NewCategoryService(categoryRepo category.RepositoryCategoryInterface) category.ServiceCategoryInterface {
	return &CategoryService{
		repo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(categoryData *entities.CategoryModels) (*entities.CategoryModels, error) {
	value := &entities.CategoryModels{
		Name:  categoryData.Name,
		Photo: categoryData.Photo,
	}

	createdCategory, err := s.repo.CreateCategory(value)
	if err != nil {
		return nil, err
	}

	return createdCategory, nil
}

func (s *CategoryService) GetAllCategory(page, perPage int) ([]*entities.CategoryModels, int64, error) {
	categories, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalCategoryCount()
	if err != nil {
		return nil, 0, err
	}

	return categories, totalItems, nil
}

func (s *CategoryService) GetCategoryByName(page int, perPage int, name string) ([]*entities.CategoryModels, int64, error) {
	categories, err := s.repo.FindByName(page, perPage, name)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalCategoryCountByName(name)
	if err != nil {
		return nil, 0, err
	}

	return categories, totalItems, nil
}

func (s *CategoryService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}

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

func (s *CategoryService) UpdateCategoryById(categoryID uint64, updatedCategory *entities.CategoryModels) error {
	categories, err := s.repo.GetCategoryById(categoryID)
	if err != nil {
		return errors.New("kategori tidak ditemukan")
	}

	err = s.repo.UpdateCategoryById(categories.ID, updatedCategory)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) DeleteCategoryById(categoryID uint64) error {
	categories, err := s.repo.GetCategoryById(categoryID)
	if err != nil {
		return errors.New("kategori tidak ditemukan")
	}

	err = s.repo.DeleteCategoryById(categories.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) GetCategoryById(categoryID uint64) (*entities.CategoryModels, error) {
	result, err := s.repo.GetCategoryById(categoryID)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return result, nil
}
