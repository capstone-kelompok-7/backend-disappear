package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.RepositoryProductInterface {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) FindByName(page, perPage int, name string) ([]entities.ProductModels, error) {
	var products []entities.ProductModels
	offset := (page - 1) * perPage
	query := r.db.Offset(offset).Limit(perPage).Preload("Categories").Preload("ProductPhotos")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&products).Error
	if err != nil {
		return products, err
	}

	return products, nil
}

func (r *ProductRepository) GetTotalProductCountByName(name string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.ProductModels{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *ProductRepository) FindAll(page, perPage int) ([]entities.ProductModels, error) {
	var products []entities.ProductModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Preload("Categories").Preload("ProductPhotos").Find(&products).Error
	if err != nil {
		return products, err
	}
	return products, nil
}

func (r *ProductRepository) GetTotalProductCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.ProductModels{}).Count(&count).Error
	return count, err
}
