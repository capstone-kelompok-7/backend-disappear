package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"gorm.io/gorm"
	"time"
)

type CarouselRepository struct {
	db *gorm.DB
}

func NewCarouselRepository(db *gorm.DB) carousel.RepositoryCarouselInterface {
	return &CarouselRepository{
		db: db,
	}
}

func (r *CarouselRepository) FindByName(page, perPage int, name string) ([]*entities.CarouselModels, error) {
	var carousels []*entities.CarouselModels
	offset := (page - 1) * perPage
	query := r.db.Offset(offset).Limit(perPage).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&carousels).Error
	if err != nil {
		return nil, err
	}

	return carousels, nil
}

func (r *CarouselRepository) GetTotalCarouselCountByName(name string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.CarouselModels{}).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *CarouselRepository) FindAll(page, perPage int) ([]*entities.CarouselModels, error) {
	var carousels []*entities.CarouselModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Where("deleted_at IS NULL").Find(&carousels).Error
	if err != nil {
		return nil, err
	}
	return carousels, nil
}

func (r *CarouselRepository) GetTotalCarouselCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.CarouselModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *CarouselRepository) CreateCarousel(carousel *entities.CarouselModels) (*entities.CarouselModels, error) {
	err := r.db.Create(&carousel).Error
	if err != nil {
		return nil, err
	}
	return carousel, nil
}

func (r *CarouselRepository) GetCarouselById(carouselID uint64) (*entities.CarouselModels, error) {
	var carousels *entities.CarouselModels
	if err := r.db.Where("id = ? AND deleted_at IS NULL", carouselID).First(&carousels).Error; err != nil {
		return nil, err
	}
	return carousels, nil
}

func (r *CarouselRepository) UpdateCarousel(carouselID uint64, updatedCarousel *entities.CarouselModels) error {
	var carousel *entities.CarouselModels
	if err := r.db.Where("id = ? AND deleted_at IS NULL", carouselID).First(&carousel).Error; err != nil {
		return err
	}
	if err := r.db.Updates(&updatedCarousel).Error; err != nil {
		return err
	}
	return nil
}

func (r *CarouselRepository) DeleteCarousel(carouselID uint64) error {
	carousels := &entities.CarouselModels{}
	if err := r.db.First(carousels, carouselID).Error; err != nil {
		return err
	}

	if err := r.db.Model(carousels).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}
