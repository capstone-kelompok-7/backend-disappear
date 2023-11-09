package category

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/category/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryCategoryInterface interface {
	CreateCategory(category *domain.CategoryModels) (*domain.CategoryModels, error)
	GetCategoryByName(name string) ([]*domain.CategoryModels, error)
	GetCategoryById(id int) (*domain.CategoryModels, error)
	FindByName(page, perPage int, name string) ([]*domain.CategoryModels, error)
	GetTotalCategoryCountByName(name string) (int64, error)
	FindAll(page, perPage int) ([]*domain.CategoryModels, error)
	GetTotalCategoryCount() (int64, error)
	UpdateCategoryById(id int, updatedCategory *domain.CategoryModels) (*domain.CategoryModels, error)
	DeleteCategoryById(id int) error
}

type ServiceCategoryInterface interface {
	CreateCategory(categoryData *domain.CategoryModels) (*domain.CategoryModels, error)
	GetAllCategory(page int, perPage int) ([]*domain.CategoryModels, int64, error)
	GetCategoryByName(page, perPage int, name string) ([]*domain.CategoryModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage int, totalPages int) int
	GetPrevPage(currentPage int) int
	UpdateCategoryById(id int, updatedCategory *domain.CategoryModels) (*domain.CategoryModels, error)
	DeleteCategoryById(id int) error
}

type HandlerCategoryInterface interface {
	CreateCategory() echo.HandlerFunc
	GetAllCategory() echo.HandlerFunc
	GetCategoryByName() echo.HandlerFunc
	UpdateCategoryById() echo.HandlerFunc
	DeleteCategoryById() echo.HandlerFunc
}
