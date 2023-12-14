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
	CreateProduct(productData *entities.ProductModels, categoryIDs []uint64) (*entities.ProductModels, error)
	GetProductByID(productID uint64) (*entities.ProductModels, error)
	CreateImageProduct(productImage *entities.ProductPhotosModels) (*entities.ProductPhotosModels, error)
	UpdateTotalReview(productID uint64) error
	UpdateProductRating(productID uint64, newRating float64) error
	GetProductReviews(page, perPage int) ([]*entities.ProductModels, error)
	UpdateProduct(product *entities.ProductModels) (*entities.ProductModels, error)
	UpdateProductCategories(product *entities.ProductModels, categoryIDs []uint64) error
	DeleteProduct(id uint64) error
	DeleteProductImage(productID, imageID uint64) error
	ReduceStockWhenPurchasing(productID, stock uint64) error
	IncreaseStock(productID, quantity uint64) error
	GetTotalProductSold() (uint64, error)
	GetTopRatedProducts() ([]*entities.ProductModels, error)
	GetProductsByCategoryAndName(page, perPage int, categoryName, name string) ([]*entities.ProductModels, error)
	GetProductsCountByCategoryAndName(categoryName, name string) (int64, error)
	GetProductsByCategoryName(categoryName string, page, perPage int) ([]*entities.ProductModels, error)
	GetProductCountByCategoryName(categoryName string) (int64, error)
	GetProductBySearchAndFilter(page, perPage int, sortBy, search string) ([]*entities.ProductModels, int64, error)
	GetProductByFilter(page, perPage int, sortBy string) ([]*entities.ProductModels, int64, error)
	GetRatedProductsInRange(page, perPage int, lowerBound, upperBound float64) ([]*entities.ProductModels, int64, error)
	SearchByNameAndFilterByRating(page, perPage int, name, ratingParam string, lowerBound, upperBound float64) ([]*entities.ProductModels, int64, error)
	FindAllProductByUserPreference(page, perPage int, productsFromAI []string) ([]*entities.ProductModels, error)
}

type ServiceProductInterface interface {
	GetAll(page, perPage int) ([]*entities.ProductModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetProductsByName(page, perPage int, name string) ([]*entities.ProductModels, int64, error)
	CreateProduct(request *dto.CreateProductRequest) (*entities.ProductModels, error)
	GetProductByID(productID uint64) (*entities.ProductModels, error)
	CreateImageProduct(request dto.CreateProductImage) (*entities.ProductPhotosModels, error)
	UpdateTotalReview(productID uint64) error
	UpdateProductRating(productID uint64, newRating float64) error
	GetProductReviews(page, perPage int) ([]*entities.ProductModels, int64, error)
	UpdateProduct(productID uint64, request *dto.UpdateProduct) (*entities.ProductModels, error)
	DeleteProduct(id uint64) error
	DeleteImageProduct(productId, imageId uint64) error
	ReduceStockWhenPurchasing(productID, quantity uint64) error
	IncreaseStock(productID, quantity uint64) error
	GetTotalProductSold() (uint64, error)
	GetTopRatedProducts() ([]*entities.ProductModels, error)
	GetProductsByCategoryAndName(categoryName, name string, page, perPage int) ([]*entities.ProductModels, int64, error)
	GetProductsByCategoryName(categoryName string, page, perPage int) ([]*entities.ProductModels, int64, error)
	GetProductsBySearchAndFilter(page, perPage int, filter, search string) ([]*entities.ProductModels, int64, error)
	GetProductsByFilter(page, perPage int, filter string) ([]*entities.ProductModels, int64, error)
	GetRatedProductsInRange(page, perPage int, ratingParam string) ([]*entities.ProductModels, int64, error)
	SearchByNameAndFilterByRating(page, perPage int, name, ratingParam string) ([]*entities.ProductModels, int64, error)
	GetProductRecommendation(userID uint64, page, perPage int) ([]*entities.ProductModels, int64, error)
}

type HandlerProductInterface interface {
	GetAllProducts() echo.HandlerFunc
	CreateProduct() echo.HandlerFunc
	GetProductById() echo.HandlerFunc
	CreateProductImage() echo.HandlerFunc
	GetAllProductsReview() echo.HandlerFunc
	UpdateProduct() echo.HandlerFunc
	DeleteProduct() echo.HandlerFunc
	DeleteProductImageById() echo.HandlerFunc
	GetAllProductsPreferences() echo.HandlerFunc
	GetTopRatedProducts() echo.HandlerFunc
}
