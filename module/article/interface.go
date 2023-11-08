package article

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/article/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryArticleInterface interface {
	FindAll(page, perpage int) ([]domain.Articles, error)
	GetTotalArticleCount() (int64, error)
	FindByTitle(page, perpage int, title string) ([]domain.Articles, error)
	GetTotalArticleCountByTitle(title string) (int64, error)
}

type ServiceArticleInterface interface {
	GetAll(page, perPage int) ([]domain.Articles, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetArticlesByTitle(page, perPage int, title string) ([]domain.Articles, int64, error)
}

type HandlerArticleInterface interface {
	GetAllArticles() echo.HandlerFunc
}
