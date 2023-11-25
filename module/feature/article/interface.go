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
	FindAll(page, perpage int) ([]entities.ArticleModels, error)
	GetTotalArticleCount() (int64, error)
	FindByTitle(page, perpage int, title string) ([]entities.ArticleModels, error)
	GetTotalArticleCountByTitle(title string) (int64, error)
	GetArticleById(id uint64) (*entities.ArticleModels, error)
	GetArticlesByDateRange(page, perpage int, startDate, endDate time.Time) ([]entities.ArticleModels, error)
	GetTotalArticleCountByDateRange(startDate, endDate time.Time) (int64, error)
}

type ServiceArticleInterface interface {
	CreateArticle(articleData *entities.ArticleModels) (*entities.ArticleModels, error)
	UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error)
	DeleteArticleById(id uint64) error
	GetAll(page, perPage int) ([]entities.ArticleModels, int64, error)
	GetArticlesByTitle(page, perPage int, title string) ([]entities.ArticleModels, int64, error)
	GetArticleById(id uint64, incrementVIews bool) (*entities.ArticleModels, error)
	GetArticlesByDateRange(page, perPage int, filterType string) ([]entities.ArticleModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
}

type HandlerArticleInterface interface {
	CreateArticle() echo.HandlerFunc
	UpdateArticleById() echo.HandlerFunc
	DeleteArticleById() echo.HandlerFunc
	GetAllArticles() echo.HandlerFunc
	GetArticleById() echo.HandlerFunc
}
