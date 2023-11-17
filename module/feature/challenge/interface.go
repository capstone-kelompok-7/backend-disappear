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
	CreateChallenge(challenge entities.ChallengeModels) (entities.ChallengeModels, error)
	UpdateChallenge(id uint64, updatedChallenge entities.ChallengeModels) (entities.ChallengeModels, error)
	GetChallengeById(id uint64) (entities.ChallengeModels, error)
	DeleteChallenge(id uint64) error
}

type ServiceChallengeInterface interface {
	GetAllChallenges(page, perPage int) ([]entities.ChallengeModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetChallengeByTitle(page, perPage int, title string) ([]entities.ChallengeModels, int64, error)
	CreateChallenge(newData entities.ChallengeModels) (entities.ChallengeModels, error)
	UpdateChallenge(id uint64, updatedChallenge entities.ChallengeModels) (entities.ChallengeModels, error)
	GetChallengeById(id uint64) (entities.ChallengeModels, error)
	DeleteChallenge(id uint64) error
}

type HandlerChallengeInterface interface {
	GetAllChallenges() echo.HandlerFunc
	CreateChallenge() echo.HandlerFunc
	UpdateChallenge() echo.HandlerFunc
	DeleteChallengeById() echo.HandlerFunc
	GetChallengeById() echo.HandlerFunc
}
