package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/homepage"
	"gorm.io/gorm"
)

type HomepageRepository struct {
	db *gorm.DB
}

func NewHomepageRepository(db *gorm.DB) homepage.RepositoryHomepageInterface {
	return &HomepageRepository{
		db: db,
	}
}

func (r *HomepageRepository) GetBestSellingProducts(limit int) ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels

	err := r.db.
		Table("products").
		Select("products.*, SUM(order_details.quantity) as total_sold").
		Joins("JOIN order_details ON order_details.product_id = products.id").
		Where("products.deleted_at IS NULL").
		Group("products.id").
		Order("total_sold desc").
		Limit(limit).
		Preload("ProductPhotos").
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *HomepageRepository) GetFiveCategories() ([]*entities.CategoryModels, error) {
	var categories []*entities.CategoryModels
	err := r.db.Where("deleted_at IS NULL").Limit(5).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *HomepageRepository) GetFiveCarousel() ([]*entities.CarouselModels, error) {
	var carousels []*entities.CarouselModels
	err := r.db.Where("deleted_at IS NULL").Limit(5).Find(&carousels).Error
	if err != nil {
		return nil, err
	}
	return carousels, nil
}

func (r *HomepageRepository) GetFiveChallenge() ([]*entities.ChallengeModels, error) {
	var challenges []*entities.ChallengeModels
	err := r.db.Where("deleted_at IS NULL").Limit(5).Find(&challenges).Error
	if err != nil {
		return nil, err
	}
	return challenges, nil
}

func (r *HomepageRepository) GetThreeArticle() ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	err := r.db.Where("deleted_at IS NULL").Limit(3).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}
