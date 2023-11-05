package auth

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryAuthInterface interface {
	Register(newData *domain.UserModels) (*domain.UserModels, error)
	Login(email string) (*domain.UserModels, error)
}

type ServiceAuthInterface interface {
	Register(newData *domain.UserModels) (*domain.UserModels, error)
	Login(email, password string) (*domain.UserModels, string, error)
}

type HandlerAuthInterface interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
}
