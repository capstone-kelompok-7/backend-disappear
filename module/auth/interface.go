package auth

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryAuthInterface interface {
	Register(newData *domain.UserModels) (*domain.UserModels, error)
	Login(email string) (*domain.UserModels, error)
	SaveOTP(otp *domain.OTPModels) (*domain.OTPModels, error)
	UpdateUser(user *domain.UserModels) (*domain.UserModels, error)
	FindValidOTP(userID int, otp string) (*domain.OTPModels, error)
	DeleteOTP(otp *domain.OTPModels) error
}

type ServiceAuthInterface interface {
	Register(newData *domain.UserModels) (*domain.UserModels, error)
	Login(email, password string) (*domain.UserModels, string, error)
	VerifyEmail(email, otp string) error
}

type HandlerAuthInterface interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	VerifyEmail() echo.HandlerFunc
}
