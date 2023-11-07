package users

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryUserInterface interface {
	GetAllUsers() ([]*domain.UserModels, error)
	GetUsersByEmail(email string) (*domain.UserModels, error)
	GetUsersById(userId uint64) (*domain.UserModels, error)
	ChangePassword(password string) (*domain.UserModels, error)
	ComparePassword(oldPass string) (*domain.UserModels, error)
}

type ServiceUserInterface interface {
	GetAllUsers() ([]*domain.UserModels, error)
	GetUsersByEmail(email string) (*domain.UserModels, error)
	GetUsersById(userId uint64) (*domain.UserModels, error)
	ChangePassword(email, oldPass, newPass string) (*domain.UserModels, error)
}

type HandlerUserInterface interface {
	GetAllUsers() echo.HandlerFunc
	GetUsersByEmail() echo.HandlerFunc
	ChangePassword() echo.HandlerFunc
}
