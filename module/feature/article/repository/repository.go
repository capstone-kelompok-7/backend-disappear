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

func (r *ArticleRepository) CreateArticle(article *entities.ArticleModels) (*entities.ArticleModels, error) {
	if err := r.db.Create(article).Error; err != nil {
		return nil, err
	}

	return article, nil
}

func (r *ArticleRepository) GetArticleById(id uint64) (*entities.ArticleModels, error) {
	var article entities.ArticleModels
	if err := r.db.Where("id =? AND deleted_at IS NULL", id).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) FindAll(page, perpage int) ([]entities.ArticleModels, error) {
	var articles []entities.ArticleModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

func (r *ArticleRepository) GetTotalArticleCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.ArticleModels{}).Count(&count).Error
	return count, err
}

func (r *ArticleRepository) FindByTitle(page, perpage int, title string) ([]entities.ArticleModels, error) {
	var articles []entities.ArticleModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("title LIKE?", "%"+title+"%").Find(&articles).Error
	if err != nil {
		return articles, err
	}
	return articles, nil
}

func (r *ArticleRepository) GetTotalArticleCountByTitle(title string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ArticleModels{}).Where("title LIKE?", "%"+title+"%").Count(&count).Error
	return count, err
}

func (r *ArticleRepository) UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error) {
	var article entities.ArticleModels
	if err := r.db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&article).Updates(updatedArticle).Error; err != nil {
		return nil, err
	}

	return updatedArticle, nil
}

func (r *ArticleRepository) DeleteArticleById(id uint64) error {
	var article entities.ArticleModels
	if err := r.db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	if err := r.db.Delete(&article).Error; err != nil {
		return err
	}

	return nil
}
