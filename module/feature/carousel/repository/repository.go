package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"gorm.io/gorm"
)

type CarouselRepository struct {
	db *gorm.DB
}

func NewCarouselRepository(db *gorm.DB) carousel.RepositoryCarouselInterface {
	return &CarouselRepository{
		db: db,
	}
}

func (r *CarouselRepository) FindByName(page, perPage int, name string) ([]entities.CarouselModels, error) {
	var carousels []entities.CarouselModels
	offset := (page - 1) * perPage
	query := r.db.Offset(offset).Limit(perPage)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&carousels).Error
	if err != nil {
		return carousels, err
	}

	return carousels, nil
}

func (r *CarouselRepository) GetTotalCarouselCountByName(name string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.CarouselModels{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *CarouselRepository) FindAll(page, perPage int) ([]entities.CarouselModels, error) {
	var carousels []entities.CarouselModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Find(&carousels).Error
	if err != nil {
		return carousels, err
	}
	return carousels, nil
}

func (r *CarouselRepository) GetTotalCarouselCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.CarouselModels{}).Count(&count).Error
	return count, err
}
