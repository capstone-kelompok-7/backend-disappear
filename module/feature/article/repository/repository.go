package repository

import (
	"time"

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

func (r *ArticleRepository) UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error) {
	var articles entities.ArticleModels
	if err := r.db.First(&articles, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&articles).Updates(updatedArticle).Error; err != nil {
		return nil, err
	}

	return updatedArticle, nil
}

func (r *ArticleRepository) UpdateArticleViews(article *entities.ArticleModels) error {
	return r.db.Save(article).Error
}

func (r *ArticleRepository) DeleteArticleById(id uint64) error {
	articles := &entities.ArticleModels{}
	if err := r.db.First(articles, id).Error; err != nil {
		return err
	}

	if err := r.db.Model(articles).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}

func (r *ArticleRepository) FindAll() ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	err := r.db.Where("deleted_at IS NULL").Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) FindByTitle(title string) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	err := r.db.Where("deleted_at IS NULL AND title LIKE?", "%"+title+"%").Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) GetArticleById(id uint64) (*entities.ArticleModels, error) {
	var articles entities.ArticleModels
	if err := r.db.Where("id =? AND deleted_at IS NULL", id).First(&articles).Error; err != nil {
		return nil, err
	}
	return &articles, nil
}

func (r *ArticleRepository) GetArticlesByDateRange(startDate, endDate time.Time) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	if err := r.db.Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *ArticleRepository) IsArticleAlreadyBookmarked(userID uint64, articleID uint64) (bool, error) {
	var exitingBookmark entities.ArticleBookmarkModels
	err := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).First(&exitingBookmark).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *ArticleRepository) BookmarkArticle(bookmarkArticle *entities.ArticleBookmarkModels) error {
	if err := r.db.Create(&bookmarkArticle).Error; err != nil {
		return err
	}
	return nil
}

func (r *ArticleRepository) DeleteBookmarkArticle(userID, articleID uint64) error {
	if err := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&entities.ArticleBookmarkModels{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *ArticleRepository) GetUserBookmarkArticle(userID uint64) ([]*entities.ArticleBookmarkModels, error) {
	var userBookmarks []*entities.ArticleBookmarkModels

	if err := r.db.Preload("Article").Where("user_id =?", userID).Find(&userBookmarks).Error; err != nil {
		return nil, err
	}
	return userBookmarks, nil
}

func (r *ArticleRepository) GetLatestArticle() ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels

	if err := r.db.Limit(5).Order("created_at desc").Where("deleted_at IS NULL").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) GetOldestArticle(page, perPage int) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	offset := (page - 1) * perPage
	if err := r.db.Offset(offset).Limit(perPage).Order("created_at asc").Where("deleted_at IS NULL").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) GetTotalArticleCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.ArticleModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *ArticleRepository) GetArticleAlphabet(page, perPage int) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels

	offset := (page - 1) * perPage

	if err := r.db.Order("title asc").Where("deleted_at IS NULL").Offset(offset).Limit(perPage).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) GetArticleMostViews(page, perPage int) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	offset := (page - 1) * perPage
	if err := r.db.Order("views desc").Where("deleted_at IS NULL").Offset(offset).Limit(perPage).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) GetOtherArticle() ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels

	if err := r.db.Order("views desc").Where("deleted_at IS NULL").Limit(5).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) SearchArticlesWithDateFilter(searchText string, startDate, endDate time.Time) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels

	query := r.db.Where("deleted_at IS NULL")

	if searchText != "" {
		query = query.Where("title LIKE ?", "%"+searchText+"%")
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("created_at <= ?", endDate)
	} else {
		query = query.Where("deleted_at IS NULL")
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) FindAllArticle(page, perPage int) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	offset := (page - 1) * perPage

	err := r.db.Where("deleted_at IS NULL").Offset(offset).Limit(perPage).Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return articles, nil
}
