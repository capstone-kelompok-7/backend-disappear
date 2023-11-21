package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category"
	"gorm.io/gorm"
	"time"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) category.RepositoryCategoryInterface {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) CreateCategory(category *entities.CategoryModels) (*entities.CategoryModels, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) GetCategoryByName(name string) ([]*entities.CategoryModels, error) {
	var categories []*entities.CategoryModels

	err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryById(id uint64) (*entities.CategoryModels, error) {
	var category entities.CategoryModels
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) UpdateCategoryById(id uint64, updatedCategory *entities.CategoryModels) (*entities.CategoryModels, error) {
	var category entities.CategoryModels
	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&category).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (r *CategoryRepository) DeleteCategoryById(id uint64) error {
	categories := &entities.CategoryModels{}
	if err := r.db.First(&categories, id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&categories).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) FindAll(page, perPage int) ([]*entities.CategoryModels, error) {
	var categories []*entities.CategoryModels
	offset := (page - 1) * perPage
	err := r.db.Limit(perPage).Offset(offset).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) FindByName(page, perPage int, name string) ([]*entities.CategoryModels, error) {
	var categories []*entities.CategoryModels
	offset := (page - 1) * perPage
	query := r.db.Offset(offset).Limit(perPage)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetTotalCategoryCountByName(name string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.CategoryModels{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *CategoryRepository) GetTotalCategoryCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.CategoryModels{}).Count(&count).Error
	return count, err
}
