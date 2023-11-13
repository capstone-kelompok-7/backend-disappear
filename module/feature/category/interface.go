package category

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryCategoryInterface interface {
	CreateCategory(category *entities.CategoryModels) (*entities.CategoryModels, error)
	GetCategoryByName(name string) ([]*entities.CategoryModels, error)
	GetCategoryById(id uint64) (*entities.CategoryModels, error)
	FindByName(page, perPage int, name string) ([]*entities.CategoryModels, error)
	GetTotalCategoryCountByName(name string) (int64, error)
	FindAll(page, perPage int) ([]*entities.CategoryModels, error)
	GetTotalCategoryCount() (int64, error)
	UpdateCategoryById(id uint64, updatedCategory *entities.CategoryModels) (*entities.CategoryModels, error)
	DeleteCategoryById(id uint64) error
}

type ServiceCategoryInterface interface {
	CreateCategory(categoryData *entities.CategoryModels) (*entities.CategoryModels, error)
	GetAllCategory(page int, perPage int) ([]*entities.CategoryModels, int64, error)
	GetCategoryByName(page, perPage int, name string) ([]*entities.CategoryModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage int, totalPages int) int
	GetPrevPage(currentPage int) int
	UpdateCategoryById(id uint64, updatedCategory *entities.CategoryModels) (*entities.CategoryModels, error)
	DeleteCategoryById(id uint64) error
	GetCategoryById(id uint64) (*entities.CategoryModels, error)
}

type HandlerCategoryInterface interface {
	CreateCategory() echo.HandlerFunc
	GetAllCategory() echo.HandlerFunc
	GetCategoryByName() echo.HandlerFunc
	UpdateCategoryById() echo.HandlerFunc
	DeleteCategoryById() echo.HandlerFunc
}
