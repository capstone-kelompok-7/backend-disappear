package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/dto"
	"math"
	"time"
)

type ProductService struct {
	repo product.RepositoryProductInterface
}

func NewProductService(repo product.RepositoryProductInterface) product.ServiceProductInterface {
	return &ProductService{
		repo: repo,
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

func (s *ProductService) CreateProduct(request *dto.CreateProductRequest) error {

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
	err := s.repo.CreateProduct(productData, request.Categories)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) GetProductByID(productID uint64) (*entities.ProductModels, error) {
	products, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return products, nil
}

func (s *ProductService) CreateImageProduct(request dto.CreateProductImage) (*entities.ProductPhotosModels, error) {
	value := &entities.ProductPhotosModels{
		ProductID: request.ProductID,
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

func (s *ProductService) UpdateProduct(productID uint64, request *dto.UpdateProduct) error {
	productData, err := s.repo.GetProductByID(productID)
	if err != nil {
		return errors.New("product tidak ditemukan")
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
		return err
	}

	err = s.repo.UpdateProduct(productData)
	if err != nil {
		return err
	}

	return nil
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
