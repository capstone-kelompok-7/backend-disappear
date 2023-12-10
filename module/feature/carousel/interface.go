package carousel

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryCarouselInterface interface {
	GetCarouselById(id uint64) (*entities.CarouselModels, error)
	FindByName(page, perPage int, name string) ([]*entities.CarouselModels, error)
	GetTotalCarouselCountByName(name string) (int64, error)
	FindAll(page, perPage int) ([]*entities.CarouselModels, error)
	GetTotalCarouselCount() (int64, error)
	CreateCarousel(carousel *entities.CarouselModels) (*entities.CarouselModels, error)
	UpdateCarousel(carouselID uint64, updatedCarousel *entities.CarouselModels) error
	DeleteCarousel(carouselID uint64) error
}

type ServiceCarouselInterface interface {
	GetAll(page, perPage int) ([]*entities.CarouselModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetCarouselsByName(page, perPage int, name string) ([]*entities.CarouselModels, int64, error)
	CreateCarousel(carouselData *entities.CarouselModels) (*entities.CarouselModels, error)
	GetCarouselById(carouselID uint64) (*entities.CarouselModels, error)
	UpdateCarousel(carouselID uint64, updatedCarousel *entities.CarouselModels) error
	DeleteCarousel(carouselID uint64) error
}
type HandlerCarouselInterface interface {
	GetAllCarousels() echo.HandlerFunc
	GetCarouselById() echo.HandlerFunc
	CreateCarousel() echo.HandlerFunc
	DeleteCarousel() echo.HandlerFunc
	UpdateCarousel() echo.HandlerFunc
}
