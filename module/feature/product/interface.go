package product

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/dto"
	"github.com/labstack/echo/v4"
)

type RepositoryProductInterface interface {
	FindByName(page, perPage int, name string) ([]*entities.ProductModels, error)
	GetTotalProductCountByName(name string) (int64, error)
	FindAll(page, perPage int) ([]*entities.ProductModels, error)
	GetTotalProductCount() (int64, error)
	CreateProduct(productData *entities.ProductModels, categoryIDs []uint64) error
	GetProductByID(productID uint64) (*entities.ProductModels, error)
	CreateImageProduct(productImage *entities.ProductPhotosModels) (*entities.ProductPhotosModels, error)
	UpdateTotalReview(productID uint64) error
	UpdateProductRating(productID uint64, newRating float64) error
	GetProductReviews(page, perPage int) ([]*entities.ProductModels, error)
	UpdateProduct(product *entities.ProductModels) error
	UpdateProductCategories(product *entities.ProductModels, categoryIDs []uint64) error
}

type ServiceProductInterface interface {
	GetAll(page, perPage int) ([]*entities.ProductModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetProductsByName(page, perPage int, name string) ([]*entities.ProductModels, int64, error)
	CreateProduct(request *dto.CreateProductRequest) error
	GetProductByID(productID uint64) (*entities.ProductModels, error)
	CreateImageProduct(request dto.CreateProductImage) (*entities.ProductPhotosModels, error)
	UpdateTotalReview(productID uint64) error
	UpdateProductRating(productID uint64, newRating float64) error
	GetProductReviews(page, perPage int) ([]*entities.ProductModels, int64, error)
	UpdateProduct(productID uint64, request *dto.UpdateProduct) error
}

type HandlerProductInterface interface {
	GetAllProducts() echo.HandlerFunc
	CreateProduct() echo.HandlerFunc
	GetProductById() echo.HandlerFunc
	CreateProductImage() echo.HandlerFunc
	GetAllProductsReview() echo.HandlerFunc
	UpdateProduct() echo.HandlerFunc
}
