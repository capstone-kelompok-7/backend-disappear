package auth

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryAuthInterface interface {
	Register(newData *entities.UserModels) (*entities.UserModels, error)
	Login(email string) (*entities.UserModels, error)
	SaveOTP(otp *entities.OTPModels) (*entities.OTPModels, error)
	UpdateUser(user *entities.UserModels) (*entities.UserModels, error)
	FindValidOTP(userID int, otp string) (*entities.OTPModels, error)
	DeleteOTP(otp *entities.OTPModels) error
	DeleteUserOTP(userId uint64) error
	ResetPassword(email, newPasswordHash string) error
}

type ServiceAuthInterface interface {
	Register(newData *entities.UserModels) (*entities.UserModels, error)
	Login(email, password string) (*entities.UserModels, string, error)
	VerifyEmail(email, otp string) error
	ResendOTP(email string) (*entities.OTPModels, error)
	ResetPassword(email, newPassword, confirmPass string) error
}

type HandlerAuthInterface interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	VerifyEmail() echo.HandlerFunc
	ResendOTP() echo.HandlerFunc
	VerifyOTP() echo.HandlerFunc
	ForgotPassword() echo.HandlerFunc
	ResetPassword() echo.HandlerFunc
}
