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
	FindAll() ([]*entities.ArticleModels, error)
	FindByTitle(title string) ([]*entities.ArticleModels, error)
	GetArticleById(id uint64) (*entities.ArticleModels, error)
	GetArticlesByDateRange(startDate, endDate time.Time) ([]*entities.ArticleModels, error)
	IsArticleAlreadyBookmarked(userID uint64, articleID uint64) (bool, error)
	BookmarkArticle(bookmarkArticle *entities.ArticleBookmarkModels) error
	DeleteBookmarkArticle(userID, articleID uint64) error
	GetUserBookmarkArticle(userID uint64) ([]*entities.ArticleBookmarkModels, error)
	GetLatestArticle() ([]*entities.ArticleModels, error)
	GetOldestArticle(page, perPage int) ([]*entities.ArticleModels, error)
	GetTotalArticleCount() (int64, error)
	GetArticleAlphabet(page, perPage int) ([]*entities.ArticleModels, error)
	GetArticleMostViews(page, perPage int) ([]*entities.ArticleModels, error)
	GetOtherArticle() ([]*entities.ArticleModels, error)
	SearchArticlesWithDateFilter(searchText string, startDate, endDate time.Time) ([]*entities.ArticleModels, error)
	FindAllArticle(page, perPage int) ([]*entities.ArticleModels, error)
}

type ServiceArticleInterface interface {
	/*?*/CreateArticle(articleData *entities.ArticleModels) (*entities.ArticleModels, error)
	/*?*/UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error)
	/*?*/DeleteArticleById(id uint64) error
	/*?*/GetAll() ([]*entities.ArticleModels, error)
	/*?*/GetArticlesByTitle(title string) ([]*entities.ArticleModels, error)
	/*?*/GetArticleById(id uint64, incrementVIews bool) (*entities.ArticleModels, error)
	/*?*/GetArticlesByDateRange(filterType string) ([]*entities.ArticleModels, error)
	/*?*/GetLatestArticles() ([]*entities.ArticleModels, error)
	/*?*/GetOldestArticle(page, perPage int) ([]*entities.ArticleModels, int64, error)
	/*?*/GetArticlesAlphabet(page, perPage int) ([]*entities.ArticleModels, int64, error)
	/*?*/GetArticleMostViews(page, perPage int) ([]*entities.ArticleModels, int64, error)
	/*?*/GetOtherArticle() ([]*entities.ArticleModels, error)
	/*1 case lagi*/GetArticleSearchByDateRange(filterType, searchText string) ([]*entities.ArticleModels, error)
	/*?*/GetAllArticleUser(page, perPage int) ([]*entities.ArticleModels, int64, error)
	/*?*/BookmarkArticle(bookmark *entities.ArticleBookmarkModels) error
	/*?*/DeleteBookmarkArticle(userID, articleID uint64) error
	/*?*/GetUserBookmarkArticle(userID uint64) ([]*entities.ArticleBookmarkModels, error)
	/*?*/GetNextPage(currentPage int, totalPages int) int
	/*?*/GetPrevPage(currentPage int) int
	/*?*/CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetFilterDateRange(filterType string) (time.Time, time.Time, error)
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
	GetOtherArticle() echo.HandlerFunc
	GetLatestArticle() echo.HandlerFunc
	GetAllArticleUser() echo.HandlerFunc
}
