package article

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryArticleInterface interface {
	CreateArticle(article *entities.ArticleModels) (*entities.ArticleModels, error)
	GetArticleById(id uint64) (*entities.ArticleModels, error)
	FindAll(page, perpage int) ([]entities.ArticleModels, error)
	GetTotalArticleCount() (int64, error)
	FindByTitle(page, perpage int, title string) ([]entities.ArticleModels, error)
	GetTotalArticleCountByTitle(title string) (int64, error)
	UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error)
	DeleteArticleById(id uint64) error
}

type ServiceArticleInterface interface {
	CreateArticle(articleData *entities.ArticleModels) (*entities.ArticleModels, error)
	GetAll(page, perPage int) ([]entities.ArticleModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetArticlesByTitle(page, perPage int, title string) ([]entities.ArticleModels, int64, error)
	GetArticleById(id uint64) (*entities.ArticleModels, error)
	UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error)
	DeleteArticleById(id uint64) error
}

type HandlerArticleInterface interface {
	CreateArticle() echo.HandlerFunc
	GetAllArticles() echo.HandlerFunc
	UpdateArticleById() echo.HandlerFunc
	DeleteArticleById() echo.HandlerFunc
}
