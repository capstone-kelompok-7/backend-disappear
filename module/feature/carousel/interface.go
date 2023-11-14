package carousel

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryCarouselInterface interface {
	FindByName(page, perPage int, name string) ([]entities.CarouselModels, error)
	GetTotalCarouselCountByName(name string) (int64, error)
	FindAll(page, perPage int) ([]entities.CarouselModels, error)
	GetTotalCarouselCount() (int64, error)
	CreateCarousel(carousel entities.CarouselModels) (entities.CarouselModels, error)
	UpdateCarousel(id uint64, updatedCarousel entities.CarouselModels) (entities.CarouselModels, error)
	DeleteCarousel(id uint64) error
	GetCarouselById(id uint64) (entities.CarouselModels, error)
}

type ServiceCarouselInterface interface {
	GetAll(page, perPage int) ([]entities.CarouselModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetCarouselsByName(page, perPage int, name string) ([]entities.CarouselModels, int64, error)
	CreateCarousel(carouselData entities.CarouselModels) (entities.CarouselModels, error)
	UpdateCarousel(id uint64, updatedCarousel entities.CarouselModels) (entities.CarouselModels, error)
	DeleteCarousel(id uint64) error
}
type HandlerCarouselInterface interface {
	GetAllCarousels() echo.HandlerFunc
	CreateCarousel() echo.HandlerFunc
	DeleteCarousel() echo.HandlerFunc
	UpdateCarousel() echo.HandlerFunc
}
