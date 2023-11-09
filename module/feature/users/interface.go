package users

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users/dto"
	"github.com/labstack/echo/v4"
)

type RepositoryUserInterface interface {
	GetAllUsers() ([]*entities.UserModels, error)
	GetUsersByEmail(email string) (*entities.UserModels, error)
	GetUsersById(userId uint64) (*entities.UserModels, error)
	GetUsersPassword(userID uint64) (string, error)
	ChangePassword(userID uint64, newPasswordHash string) error
}

type ServiceUserInterface interface {
	GetAllUsers() ([]*entities.UserModels, error)
	GetUsersByEmail(email string) (*entities.UserModels, error)
	GetUsersById(userId uint64) (*entities.UserModels, error)
	ValidatePassword(userID uint64, oldPassword, newPassword, confirmPassword string) error
	ChangePassword(userID uint64, updateRequest dto.UpdatePasswordRequest) error
}

type HandlerUserInterface interface {
	GetAllUsers() echo.HandlerFunc
	GetUsersByEmail() echo.HandlerFunc
	ChangePassword() echo.HandlerFunc
	GetUsersById() echo.HandlerFunc
}
