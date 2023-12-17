package auth

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/dto"
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
	LoginSocial(socialID string) (*entities.UserModels, error)
	FindUserBySocialID(socialID string) (*entities.UserModels, error)
	UpdateLastLogin(userID uint64, lastLogin time.Time) error
	CekDeviceTokenByEmail(email string) (string, error)
	UpdateDeviceTokenByID(email string, deviceToken string) (*entities.UserModels, error)
}

type ServiceAuthInterface interface {
	Register(newData *entities.UserModels) (*entities.UserModels, error)
	Login(email, password, deviceToken string) (*entities.UserModels, string, error)
	VerifyEmail(email, otp string) error
	ResendOTP(email string) (*entities.OTPModels, error)
	ResetPassword(email, newPassword, confirmPass string) error
	VerifyOTP(email, otp string) (string, error)
	RegisterSocial(req *dto.RegisterSocialRequest) (*entities.UserModels, error)
	LoginSocial(socialID string) (*entities.UserModels, string, error)
}

type HandlerAuthInterface interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	VerifyEmail() echo.HandlerFunc
	ResendOTP() echo.HandlerFunc
	VerifyOTP() echo.HandlerFunc
	ForgotPassword() echo.HandlerFunc
	ResetPassword() echo.HandlerFunc
	RegisterSocial() echo.HandlerFunc
	LoginSocial() echo.HandlerFunc
}
