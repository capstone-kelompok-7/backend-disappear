package service

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"math"
)

type CarouselService struct {
	repo carousel.RepositoryCarouselInterface
}

func NewCarouselService(repo carousel.RepositoryCarouselInterface) carousel.ServiceCarouselInterface {
	return &CarouselService{
		repo: repo,
	}
}

func (s *CarouselService) GetAll(page, perPage int) ([]entities.CarouselModels, int64, error) {
	carousels, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return carousels, 0, err
	}

	totalItems, err := s.repo.GetTotalCarouselCount()
	if err != nil {
		return carousels, 0, err
	}

	return carousels, totalItems, nil
}

func (s *CarouselService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
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

func (s *CarouselService) GetNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}
	return totalPages
}

func (s *CarouselService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}
	return 1
}

func (s *CarouselService) GetCarouselsByName(page, perPage int, name string) ([]entities.CarouselModels, int64, error) {
	carousels, err := s.repo.FindByName(page, perPage, name)
	if err != nil {
		return carousels, 0, err
	}

	totalItems, err := s.repo.GetTotalCarouselCountByName(name)
	if err != nil {
		return carousels, 0, err
	}

	return carousels, totalItems, nil
}
