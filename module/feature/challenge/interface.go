package challenge

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryChallengeInterface interface {
	FindAll(page, perpage int) ([]entities.ChallengeModels, error)
	GetTotalChallengeCount() (int64, error)
	FindByTitle(page, perpage int, title string) ([]entities.ChallengeModels, error)
	GetTotalChallengeCountByTitle(title string) (int64, error)
}

type ServiceChallengeInterface interface {
	GetAllChallenges(page, perPage int) ([]entities.ChallengeModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetChallengeByTitle(page, perPage int, title string) ([]entities.ChallengeModels, int64, error)
}

type HandlerChallengeInterface interface {
	GetAllChallenges() echo.HandlerFunc
}
