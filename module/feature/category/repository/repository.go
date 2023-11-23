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

func (r *CategoryRepository) GetCategoryById(categoryID uint64) (*entities.CategoryModels, error) {
	var categories entities.CategoryModels
	if err := r.db.Where("id = ? AND deleted_at IS NULL", categoryID).First(&categories).Error; err != nil {
		return nil, err
	}
	return &categories, nil
}

func (r *CategoryRepository) UpdateCategoryById(categoryID uint64, updatedCategory *entities.CategoryModels) error {
	var categories *entities.CategoryModels
	if err := r.db.Where("id = ? AND deleted_at IS NULL", categoryID).First(&categories).Error; err != nil {
		return err
	}
	if err := r.db.Updates(&updatedCategory).Error; err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) DeleteCategoryById(categoryID uint64) error {
	categories := &entities.CategoryModels{}
	if err := r.db.First(&categories, categoryID).Error; err != nil {
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
	err := r.db.Limit(perPage).Offset(offset).Where("deleted_at IS NULL").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) FindByName(page, perPage int, name string) ([]*entities.CategoryModels, error) {
	var categories []*entities.CategoryModels
	offset := (page - 1) * perPage
	query := r.db.Offset(offset).Limit(perPage).Where("deleted_at IS NULL")

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
	query := r.db.Model(&entities.CategoryModels{}).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *CategoryRepository) GetTotalCategoryCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.CategoryModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}
