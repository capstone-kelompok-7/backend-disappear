package article

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryArticleInterface interface {
	CreateArticle(article *entities.ArticleModels) (*entities.ArticleModels, error)
	UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error)
	UpdateArticleViews(article *entities.ArticleModels) error
	DeleteArticleById(id uint64) error
	FindAll() ([]entities.ArticleModels, error)
	FindByTitle(title string) ([]entities.ArticleModels, error)
	GetArticleById(id uint64) (*entities.ArticleModels, error)
	GetArticlesByDateRange(startDate, endDate time.Time) ([]entities.ArticleModels, error)
	IsArticleAlreadyBookmarked(userID uint64, articleID uint64) (bool, error)
	BookmarkArticle(bookmarkArticle *entities.ArticleBookmarkModels) error 
	DeleteBookmarkArticle(userID, articleID uint64) error
	GetUserBookmarkArticle(userID uint64) ([]*entities.ArticleBookmarkModels, error)
	GetLatestArticle() ([]entities.ArticleModels, error) 
	GetArticleByViewsAsc() ([]entities.ArticleModels, error)
	GetArticleByViewsDesc() ([]entities.ArticleModels, error)
	GetArticleByTitleAsc() ([]entities.ArticleModels, error)
	GetArticleByTitleDesc() ([]entities.ArticleModels, error)
}

type ServiceArticleInterface interface {
	CreateArticle(articleData *entities.ArticleModels) (*entities.ArticleModels, error)
	UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error)
	DeleteArticleById(id uint64) error
	GetAll() ([]entities.ArticleModels, error)
	GetArticlesByTitle(title string) ([]entities.ArticleModels, error)
	GetArticleById(id uint64, incrementVIews bool) (*entities.ArticleModels, error)
	GetArticlesByDateRange(filterType string) ([]entities.ArticleModels, error)
	BookmarkArticle(bookmark *entities.ArticleBookmarkModels) error
	DeleteBookmarkArticle(userID, articleID uint64) error
	GetUserBookmarkArticle(userID uint64) ([]*entities.ArticleBookmarkModels, error)
	GetLatestArticles() ([]entities.ArticleModels, error)
	GetArticlesByViews(sortType string) ([]entities.ArticleModels, error)
	GetArticlesBySortedTitle(sortType string) ([]entities.ArticleModels, error)
}

type HandlerArticleInterface interface {
	CreateArticle() echo.HandlerFunc
	UpdateArticleById() echo.HandlerFunc
	DeleteArticleById() echo.HandlerFunc
	GetAllArticles() echo.HandlerFunc
	GetArticleById() echo.HandlerFunc
	BookmarkArticle() echo.HandlerFunc
	DeleteBookmarkedArticle() echo.HandlerFunc
	GetUsersBookmark() echo.HandlerFunc
}
