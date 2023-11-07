package product

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/product/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryProductInterface interface {
	FindByName(page, perPage int, name string) ([]domain.Product, error)
	GetTotalProductCountByName(name string) (int64, error)
	FindAll(page, perPage int) ([]domain.Product, error)
	GetTotalProductCount() (int64, error)
}

type ServiceProductInterface interface {
	GetAll(page, perPage int) ([]domain.Product, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetProductsByName(page, perPage int, name string) ([]domain.Product, int64, error)
}

type HandlerProductInterface interface {
	GetAllProducts() echo.HandlerFunc
}
