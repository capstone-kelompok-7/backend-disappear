package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant"
	"math"
	"strings"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/dto"
)

type ProductService struct {
	repo       product.RepositoryProductInterface
	botService assistant.ServiceAssistantInterface
}

func NewProductService(repo product.RepositoryProductInterface, botService assistant.ServiceAssistantInterface) product.ServiceProductInterface {
	return &ProductService{
		repo:       repo,
		botService: botService,
	}
}

func (s *ProductService) GetAll(page, perPage int) ([]*entities.ProductModels, int64, error) {
	products, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return products, 0, err
	}

	totalItems, err := s.repo.GetTotalProductCount()
	if err != nil {
		return products, 0, err
	}

	return products, totalItems, nil
}

func (s *ProductService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	pageInt := page
	if pageInt <= 0 {
		pageInt = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	if pageInt > totalPages {
		pageInt = totalPages
	}

	return pageInt, totalPages
}

func (s *ProductService) GetNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}
	return totalPages
}

func (s *ProductService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}
	return 1
}

func (s *ProductService) GetProductsByName(page, perPage int, name string) ([]*entities.ProductModels, int64, error) {
	products, err := s.repo.FindByName(page, perPage, name)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalProductCountByName(name)
	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (s *ProductService) CreateProduct(request *dto.CreateProductRequest) (*entities.ProductModels, error) {
	productData := &entities.ProductModels{
		Name:        request.Name,
		Description: request.Description,
		GramPlastic: request.GramPlastic,
		Price:       request.Price,
		Stock:       request.Stock,
		Discount:    request.Discount,
		Exp:         request.Exp,
		Rating:      0.0,
		TotalReview: 0,
		CreatedAt:   time.Now(),
	}
	if request.ImageURL != "" {
		productData.ProductPhotos = []entities.ProductPhotosModels{
			{
				ImageURL:  request.ImageURL,
				CreatedAt: time.Now(),
			},
		}
	}
	createdProduct, err := s.repo.CreateProduct(productData, request.Categories)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (s *ProductService) GetProductByID(productID uint64) (*entities.ProductModels, error) {
	products, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return products, nil
}

func (s *ProductService) CreateImageProduct(request dto.CreateProductImage) (*entities.ProductPhotosModels, error) {
	products, err := s.repo.GetProductByID(request.ProductID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	value := &entities.ProductPhotosModels{
		ProductID: products.ID,
		ImageURL:  request.Image,
	}

	images, err := s.repo.CreateImageProduct(value)
	if err != nil {
		return images, err
	}

	return images, nil
}

func (s *ProductService) UpdateTotalReview(productID uint64) error {
	_, err := s.repo.GetProductByID(productID)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}
	err = s.repo.UpdateTotalReview(productID)
	if err != nil {
		return errors.New("gagal memperbarui total review")
	}

	return nil
}

func (s *ProductService) UpdateProductRating(productID uint64, newRating float64) error {
	err := s.repo.UpdateProductRating(productID, newRating)
	if err != nil {
		return errors.New("gagal memperbarui rating produk")
	}

	return nil
}

func (s *ProductService) GetProductReviews(page, perPage int) ([]*entities.ProductModels, int64, error) {
	products, err := s.repo.GetProductReviews(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalProductCount()
	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (s *ProductService) UpdateProduct(productID uint64, request *dto.UpdateProduct) (*entities.ProductModels, error) {
	productData, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	productData.Name = request.Name
	productData.Description = request.Description
	productData.GramPlastic = request.GramPlastic
	productData.Price = request.Price
	productData.Stock = request.Stock
	productData.Discount = request.Discount
	productData.Exp = request.Exp
	productData.UpdatedAt = time.Now()

	if request.ImageURL != "" {
		productData.ProductPhotos = []entities.ProductPhotosModels{
			{
				ImageURL:  request.ImageURL,
				CreatedAt: time.Now(),
			},
		}
	}

	err = s.repo.UpdateProductCategories(productData, request.CategoryIDs)
	if err != nil {
		return nil, err
	}

	productData, err = s.repo.UpdateProduct(productData)
	if err != nil {
		return nil, err
	}

	return productData, nil
}

func (s *ProductService) DeleteProduct(id uint64) error {
	productId, err := s.repo.GetProductByID(id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}
	if err := s.repo.DeleteProduct(productId.ID); err != nil {
		return errors.New("gagal menghapus product")
	}
	return nil
}

func (s *ProductService) DeleteImageProduct(productId, imageId uint64) error {
	productData, err := s.repo.GetProductByID(productId)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}
	found := false
	for _, photo := range productData.ProductPhotos {
		if photo.ID == imageId {
			found = true
			break
		}
	}
	if !found {
		return errors.New("image tidak ditemukan pada produk ini")
	}

	if err := s.repo.DeleteProductImage(productData.ID, imageId); err != nil {
		return errors.New("gagal menghapus image pada produk ini")
	}
	return nil
}

func (s *ProductService) ReduceStockWhenPurchasing(productID, quantity uint64) error {
	products, err := s.repo.GetProductByID(productID)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	if products.Stock < quantity {
		return errors.New("stok tidak mencukupi untuk pesanan ini")
	}

	if err := s.repo.ReduceStockWhenPurchasing(products.ID, quantity); err != nil {
		return err
	}
	return nil
}

func (s *ProductService) IncreaseStock(productID, quantity uint64) error {
	products, err := s.repo.GetProductByID(productID)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	err = s.repo.IncreaseStock(products.ID, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) GetTotalProductSold() (uint64, error) {
	totalSold, err := s.repo.GetTotalProductSold()
	if err != nil {
		return 0, err
	}
	return totalSold, nil
}

func (s *ProductService) GetTopRatedProducts() ([]*entities.ProductModels, error) {
	result, err := s.repo.GetTopRatedProducts()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProductService) GetProductsByCategoryAndName(categoryName, name string, page, perPage int) ([]*entities.ProductModels, int64, error) {
	var products []*entities.ProductModels
	var err error

	products, err = s.repo.GetProductsByCategoryAndName(page, perPage, categoryName, name)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetProductsCountByCategoryAndName(categoryName, name)
	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (s *ProductService) GetProductsByCategoryName(categoryName string, page, perPage int) ([]*entities.ProductModels, int64, error) {
	products, err := s.repo.GetProductsByCategoryName(categoryName, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.GetProductCountByCategoryName(categoryName)
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (s *ProductService) GetProductsBySearchAndFilter(page, perPage int, filter, search string) ([]*entities.ProductModels, int64, error) {
	var products []*entities.ProductModels
	var totalItems int64
	var err error

	if !isValidFilter(filter) {
		products, totalItems, err = s.repo.GetProductBySearchAndFilter(page, perPage, filter, search)
	} else {
		products, totalItems, err = s.repo.GetProductBySearchAndFilter(page, perPage, "", search)
	}

	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (s *ProductService) GetProductsByFilter(page, perPage int, filter string) ([]*entities.ProductModels, int64, error) {
	var products []*entities.ProductModels
	var totalItems int64
	var err error

	if !isValidFilter(filter) {
		products, totalItems, err = s.repo.GetProductByFilter(page, perPage, filter)
	}
	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func isValidFilter(filter string) bool {
	validFilters := map[string]bool{
		"abjad":    true,
		"termurah": true,
		"termahal": true,
		"terbaru":  true,
		"promo":    true,
	}
	return !validFilters[filter]
}

func getRatingBounds(ratingParam string) (float64, float64, error) {
	ratingParam = strings.ToLower(ratingParam)
	var lowerBound, upperBound float64

	switch ratingParam {
	case "sangat buruk":
		lowerBound = 0.1
		upperBound = 1.9
	case "buruk":
		lowerBound = 2
		upperBound = 2.9
	case "sedang":
		lowerBound = 3
		upperBound = 3.9
	case "baik":
		lowerBound = 4
		upperBound = 4.9
	case "sangat baik":
		lowerBound = 5
		upperBound = 5.0
	default:
		return 0, 0, errors.New("tipe filter tidak valid")
	}

	return lowerBound, upperBound, nil
}

func (s *ProductService) GetRatedProductsInRange(page, perPage int, ratingParam string) ([]*entities.ProductModels, int64, error) {
	lowerBound, upperBound, err := getRatingBounds(ratingParam)
	if err != nil {
		return nil, 0, err
	}

	products, totalItems, err := s.repo.GetRatedProductsInRange(page, perPage, lowerBound, upperBound)
	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (s *ProductService) SearchByNameAndFilterByRating(page, perPage int, name, ratingParam string) ([]*entities.ProductModels, int64, error) {
	lowerBound, upperBound, err := getRatingBounds(ratingParam)
	if err != nil {
		return nil, 0, err
	}

	products, totalItems, err := s.repo.SearchByNameAndFilterByRating(page, perPage, name, ratingParam, lowerBound, upperBound)
	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (s *ProductService) GetProductRecommendation(userID uint64, page, perPage int) ([]*entities.ProductModels, int64, error) {
	recommendations, err := s.botService.GenerateRecommendationProduct(userID)
	if err != nil {
		return nil, 0, err
	}

	result, err := s.repo.FindAllProductByUserPreference(page, perPage, recommendations)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalProductCount()
	if err != nil {
		return nil, 0, err
	}

	return result, totalItems, nil
}
