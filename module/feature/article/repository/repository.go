package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/article"
	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) article.RepositoryArticleInterface {
	return &ArticleRepository{
		db: db,
	}
}

func (r *ArticleRepository) FindAll(page, perpage int) ([]entities.Articles, error) {
	var articles []entities.Articles
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

func (r *ArticleRepository) GetTotalArticleCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.Articles{}).Count(&count).Error
	return count, err
}

func (r *ArticleRepository) FindByTitle(page, perpage int, title string) ([]entities.Articles, error) {
	var articles []entities.Articles
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("title LIKE?", "%"+title+"%").Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

func (r *ArticleRepository) GetTotalArticleCountByTitle(title string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.Articles{}).Where("title LIKE?", "%"+title+"%").Count(&count).Error
	return count, err
}
