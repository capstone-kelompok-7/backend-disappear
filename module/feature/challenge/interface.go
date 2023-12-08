package challenge

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/dto"
	"github.com/labstack/echo/v4"
)

type RepositoryChallengeInterface interface {
	FindAll(page, perpage int) ([]*entities.ChallengeModels, error)
	GetTotalChallengeCount() (int64, error)
	FindByTitle(page, perpage int, title string) ([]*entities.ChallengeModels, error)
	FindByStatus(page, perpage int, status string) ([]*entities.ChallengeModels, error)
	GetTotalChallengeCountByTitle(title string) (int64, error)
	GetTotalChallengeCountByStatus(status string) (int64, error)
	CreateChallenge(challenge *entities.ChallengeModels) (*entities.ChallengeModels, error)
	UpdateChallenge(id uint64, updatedChallenge *entities.ChallengeModels) (*entities.ChallengeModels, error)
	GetChallengeById(id uint64) (*entities.ChallengeModels, error)
	DeleteChallenge(id uint64) error
	CreateSubmitChallengeForm(form *entities.ChallengeFormModels) (*entities.ChallengeFormModels, error)
	GetAllSubmitChallengeForm(page, perpage int) ([]*entities.ChallengeFormModels, error)
	GetSubmitChallengeFormByStatus(page, perpage int, status string) ([]*entities.ChallengeFormModels, error)
	GetTotalSubmitChallengeFormCount() (int64, error)
	GetTotalSubmitChallengeFormCountByStatus(status string) (int64, error)
	GetSubmitChallengeFormById(id uint64) (*entities.ChallengeFormModels, error)
	UpdateSubmitChallengeForm(id uint64, updatedStatus dto.UpdateChallengeFormStatusRequest) (*entities.ChallengeFormModels, error)
	GetSubmitChallengeFormByUserAndChallenge(userID uint64) ([]*entities.ChallengeFormModels, error)
	GetSubmitChallengeFormByDateRange(page, perpage int, startDate, endDate time.Time) ([]*entities.ChallengeFormModels, error)
	GetTotalSubmitChallengeFormCountByDateRange(startDate, endDate time.Time) (int64, error)
	GetSubmitChallengeFormByStatusAndDate(page, perPage int, filterStatus string, startDate, endDate time.Time) ([]*entities.ChallengeFormModels, error)
	GetTotalSubmitChallengeFormCountByStatusAndDate(filterStatus string, startDate, endDate time.Time) (int64, error)
	GetChallengesBySearchAndStatus(page, perPage int, search, status string) ([]*entities.ChallengeModels, int64, error)
}

type ServiceChallengeInterface interface {
	GetAllChallenges(page, perPage int) ([]*entities.ChallengeModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetChallengeByTitle(page, perPage int, title string) ([]*entities.ChallengeModels, int64, error)
	GetChallengeByStatus(page, perPage int, status string) ([]*entities.ChallengeModels, int64, error)
	CreateChallenge(newData *entities.ChallengeModels) (*entities.ChallengeModels, error)
	UpdateChallenge(id uint64, updatedChallenge *entities.ChallengeModels) (*entities.ChallengeModels, error)
	GetChallengeById(id uint64) (*entities.ChallengeModels, error)
	DeleteChallenge(id uint64) error
	CreateSubmitChallengeForm(form *entities.ChallengeFormModels) (*entities.ChallengeFormModels, error)
	GetAllSubmitChallengeForm(page, perPage int) ([]*entities.ChallengeFormModels, int64, error)
	GetSubmitChallengeFormByStatus(page, perPage int, status string) ([]*entities.ChallengeFormModels, int64, error)
	GetSubmitChallengeFormById(id uint64) (*entities.ChallengeFormModels, error)
	UpdateSubmitChallengeForm(id uint64, updatedStatus dto.UpdateChallengeFormStatusRequest) (*entities.ChallengeFormModels, error)
	GetSubmitChallengeFormByDateRange(page, perPage int, filterType string) ([]*entities.ChallengeFormModels, int64, error)
	GetSubmitChallengeFormByStatusAndDate(page, perPage int, filterStatus string, filterType string) ([]*entities.ChallengeFormModels, int64, error)
	GetChallengesBySearchAndStatus(page, perPage int, search, status string) ([]*entities.ChallengeModels, int64, error)
}

type HandlerChallengeInterface interface {
	GetAllChallenges() echo.HandlerFunc
	CreateChallenge() echo.HandlerFunc
	UpdateChallenge() echo.HandlerFunc
	DeleteChallengeById() echo.HandlerFunc
	GetChallengeById() echo.HandlerFunc
	CreateSubmitChallengeForm() echo.HandlerFunc
	GetAllSubmitChallengeForm() echo.HandlerFunc
	UpdateSubmitChallengeForm() echo.HandlerFunc
	GetSubmitChallengeFormById() echo.HandlerFunc
}
