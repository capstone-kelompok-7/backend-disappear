package service

import (
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

func (s *ProductService) GetAll(page, perPage int) ([]entities.ProductModels, int64, error) {
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

	total_pages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	if pageInt > total_pages {
		pageInt = total_pages
	}

	return pageInt, total_pages
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

func (s *ProductService) GetProductsByName(page, perPage int, name string) ([]entities.ProductModels, int64, error) {
	products, err := s.repo.FindByName(page, perPage, name)
	if err != nil {
		return products, 0, err
	}

	totalItems, err := s.repo.GetTotalProductCountByName(name)
	if err != nil {
		return products, 0, err
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
		CreatedAt:   time.Now(),
	}

	err := s.repo.CreateProduct(productData, request.Categories)
	if err != nil {
		return err
	}

	return nil
}
