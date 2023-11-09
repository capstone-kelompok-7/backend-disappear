package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/category"
	"github.com/capstone-kelompok-7/backend-disappear/module/category/domain"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) category.RepositoryCategoryInterface {
	return &CategoryRepository{
		db: db,
	}
}
func (r *CategoryRepository) CreateCategory(category *domain.CategoryModels) (*domain.CategoryModels, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) GetCategoryByName(name string) ([]*domain.CategoryModels, error) {
	var categories []*domain.CategoryModels

	err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryById(id int) (*domain.CategoryModels, error) {
	var category domain.CategoryModels
	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) UpdateCategoryById(id int, updatedCategory *domain.CategoryModels) (*domain.CategoryModels, error) {
	var category domain.CategoryModels
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

func (r *CategoryRepository) DeleteCategoryById(id int) error {
	var category domain.CategoryModels
	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	if err := r.db.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) FindAll(page, perPage int) ([]*domain.CategoryModels, error) {
	var categories []*domain.CategoryModels
	offset := (page - 1) * perPage
	err := r.db.Limit(perPage).Offset(offset).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) FindByName(page, perPage int, name string) ([]*domain.CategoryModels, error) {
	var categories []*domain.CategoryModels
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
	query := r.db.Model(&domain.CategoryModels{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}
func (r *CategoryRepository) GetTotalCategoryCount() (int64, error) {
	var count int64
	err := r.db.Model(&domain.CategoryModels{}).Count(&count).Error
	return count, err
}
