package users

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users/dto"
	"github.com/labstack/echo/v4"
)

type RepositoryUserInterface interface {
	GetUsersByEmail(email string) (*entities.UserModels, error)
	GetUsersById(userId uint64) (*entities.UserModels, error)
	GetUsersPassword(userID uint64) (string, error)
	ChangePassword(userID uint64, newPasswordHash string) error
	FindAll(page, perPage int) ([]*entities.UserModels, error)
	FindByName(page, perPage int, name string) ([]*entities.UserModels, error)
	GetTotalUserCountByName(name string) (int64, error)
	GetTotalUserCount() (int64, error)
	EditProfile(userID uint64, updatedData dto.EditProfileRequest) (*entities.UserModels, error)
	DeleteAccount(userID uint64) error
	UpdateUserExp(userID uint64, exp uint64) (*entities.UserModels, error)
	UpdateUserChallengeFollow(userID uint64, totalChallenge uint64) (*entities.UserModels, error)
	UpdateUserContribution(userID uint64, gramPlastic uint64) (*entities.UserModels, error)
	UpdateUserLevel(userID uint64, level string) error
	GetUserLevel(userID uint64) (string, error)
	GetFilterLevel(level string) ([]*entities.UserModels, error)
	GetLeaderboardByExp(limit int) ([]*entities.UserModels, error)
	GetUserTransactionActivity(userID uint64) (int, int, int, error)
	GetUserChallengeActivity(userID uint64) (int, int, int, error)
	FindAllPersonalization() ([]*entities.PersonalizationModels, error)
	FindAllEnvironmentsIsues() ([]*entities.EnvironmentIssuesModels, error)
	CreateUserPersonalization(userID uint64, req *dto.UserPersonalizationRequest) ([]*entities.PersonalizationModels, error)
	GetUserPersonalization(userID uint64) ([]*entities.PersonalizationModels, error)
}

type ServiceUserInterface interface {
	GetUsersByEmail(email string) (*entities.UserModels, error)
	GetUsersById(userId uint64) (*entities.UserModels, error)
	ValidatePassword(userID uint64, oldPassword, newPassword, confirmPassword string) error
	ChangePassword(userID uint64, updateRequest dto.UpdatePasswordRequest) error
	GetAllUsers(page, perPage int) ([]*entities.UserModels, int64, error)
	GetUsersByName(page int, perPage int, name string) ([]*entities.UserModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage int, totalPages int) int
	GetPrevPage(currentPage int) int
	EditProfile(userID uint64, updatedData dto.EditProfileRequest) (*entities.UserModels, error)
	DeleteAccount(userID uint64) error
	UpdateUserExp(userID uint64, exp uint64) (*entities.UserModels, error)
	UpdateUserChallengeFollow(userID uint64, totalChallenge uint64) (*entities.UserModels, error)
	UpdateUserContribution(userID uint64, gramPlastic uint64) (*entities.UserModels, error)
	GetUserLevel(userID uint64) (string, error)
	GetLeaderboardByExp(limit int) ([]*entities.UserModels, error)
	GetUserTransactionActivity(userID uint64) (int, int, int, error)
	GetUserChallengeActivity(userID uint64) (int, int, int, error)
	GetUserProfile(userID uint64) (*entities.UserModels, error)
	GetAllUserPersonalization() ([]*entities.PersonalizationModels, error)
	GetAllEnvironmentsIsues() ([]*entities.EnvironmentIssuesModels, error)
	CreateUserPersonalization(userID uint64, req *dto.UserPersonalizationRequest) ([]*entities.PersonalizationModels, error)
	GetUserPersonalization(userID uint64) ([]*entities.PersonalizationModels, error)
}

type HandlerUserInterface interface {
	GetAllUsers() echo.HandlerFunc
	GetUsersByEmail() echo.HandlerFunc
	ChangePassword() echo.HandlerFunc
	GetUsersById() echo.HandlerFunc
	EditProfile() echo.HandlerFunc
	DeleteAccount() echo.HandlerFunc
	GetLeaderboard() echo.HandlerFunc
	GetUserTransactionActivity() echo.HandlerFunc
	GetUserProfile() echo.HandlerFunc

	GetAllUserPersonalization() echo.HandlerFunc
	GetAllEnvironmentsIsues() echo.HandlerFunc
	CreateUserPersonalization() echo.HandlerFunc
}
