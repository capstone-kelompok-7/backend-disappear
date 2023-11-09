package users

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryUserInterface interface {
	GetAllUsers() ([]*domain.UserModels, error)
	GetUsersByEmail(email string) (*domain.UserModels, error)
	GetUsersById(userId uint64) (*domain.UserModels, error)
	GetUsersPassword(userID uint64) (string, error)
	ChangePassword(userID uint64, newPasswordHash string) error
}

type ServiceUserInterface interface {
	GetAllUsers() ([]*domain.UserModels, error)
	GetUsersByEmail(email string) (*domain.UserModels, error)
	GetUsersById(userId uint64) (*domain.UserModels, error)
	ValidatePassword(userID uint64, oldPassword, newPassword, confirmPassword string) error
	ChangePassword(userID uint64, updateRequest domain.UpdatePasswordRequest) error
}

type HandlerUserInterface interface {
	GetAllUsers() echo.HandlerFunc
	GetUsersByEmail() echo.HandlerFunc
	ChangePassword() echo.HandlerFunc
	GetUsersById() echo.HandlerFunc
}
