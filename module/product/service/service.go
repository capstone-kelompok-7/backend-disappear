package service

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/product/domain"
	"math"
)

type ProductService struct {
	repo product.RepositoryProductInterface
}

func NewProductService(repo product.RepositoryProductInterface) product.ServiceProductInterface {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetAll(page, perPage int) ([]domain.Product, int64, error) {
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

func (s *ProductService) GetProductsByName(page, perPage int, name string) ([]domain.Product, int64, error) {
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
