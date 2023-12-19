package handler

import (
	user "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth"
	dto2 "github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/email"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	service     auth.ServiceAuthInterface
	userService users.ServiceUserInterface
}

func NewAuthHandler(service auth.ServiceAuthInterface, userService users.ServiceUserInterface) auth.HandlerAuthInterface {
	return &AuthHandler{
		service:     service,
		userService: userService,
	}
}

func (h *AuthHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		registerRequest := new(dto2.RegisterRequest)
		if err := c.Bind(registerRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(registerRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		newUser := &user.UserModels{
			Email:    registerRequest.Email,
			Password: registerRequest.Password,
		}

		_, err := h.service.Register(newUser)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendaftarkan akun: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Registrasi berhasil! Silakan cek email anda untuk OTP.")
	}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginRequest dto2.LoginRequest
		if err := c.Bind(&loginRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(loginRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		userLogin, accessToken, err := h.service.Login(loginRequest.Email, loginRequest.Password, loginRequest.DeviceToken)
		if err != nil {
			if err.Error() == "user tidak ditemukan" {
				return response.SendStatusNotFoundResponse(c, "Pengguna tidak ditemukan")
			} else if err.Error() == "akun anda belum diverifikasi" {
				return response.SendStatusUnauthorizedResponse(c, "akun anda belum diverifikasi")
			}

			logrus.Error("Kesalahan : " + err.Error())
			return response.SendStatusUnauthorizedResponse(c, "Email atau kata sandi salah")
		}

		result := &dto2.LoginResponse{
			Email:       userLogin.Email,
			AccessToken: accessToken,
		}

		return response.SendSuccessResponse(c, "Selamat datang!, Anda telah berhasil masuk.", result)
	}
}

func (h *AuthHandler) VerifyEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		var emailRequest dto2.EmailRequest

		if err := c.Bind(&emailRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}
		if err := utils.ValidateStruct(emailRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		if err := h.service.VerifyEmail(emailRequest.Email, emailRequest.OTP); err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memverifikasi email: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Email berhasil diverifikasi!, Silahkan login.")
	}
}

func (h *AuthHandler) ResendOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		var emailRequest dto2.ResendOTPRequest

		if err := c.Bind(&emailRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(emailRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		newOTP, err := h.service.ResendOTP(emailRequest.Email)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mengirim ulang OTP: "+err.Error())
		}

		emailSender := email.NewEmailService()
		err = emailSender.EmailService(emailRequest.Email, newOTP.OTP)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mengirim ulang OTP ke email: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "OTP berhasil dikirim kembali!, Silahkan cek email anda")

	}
}

func (h *AuthHandler) ForgotPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var forgotPasswordRequest dto2.ForgotPasswordRequest
		if err := c.Bind(&forgotPasswordRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(forgotPasswordRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		newOTP, err := h.service.ResendOTP(forgotPasswordRequest.Email)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mengirim ulang OTP: "+err.Error())
		}

		emailSender := email.NewEmailService()
		err = emailSender.EmailService(forgotPasswordRequest.Email, newOTP.OTP)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mengirim ulang OTP ke email: "+err.Error())

		}

		return response.SendStatusOkResponse(c, "Permintaan OTP berhasil dikirim")
	}
}

func (h *AuthHandler) VerifyOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		var emailRequest dto2.EmailRequest

		if err := c.Bind(&emailRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}
		if err := utils.ValidateStruct(emailRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		accessToken, err := h.service.VerifyOTP(emailRequest.Email, emailRequest.OTP)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal verifikasi OTP: "+err.Error())
		}

		result := &dto2.VerifyOTPResponse{
			AccessToken: accessToken,
		}
		return response.SendSuccessResponse(c, "Verifikasi OTP berhasil", result)
	}
}

func (h *AuthHandler) ResetPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var resetPasswordRequest dto2.ResetPasswordRequest
		if err := c.Bind(&resetPasswordRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(resetPasswordRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		currentUser := c.Get("CurrentUser").(*user.UserModels)
		err := h.service.ResetPassword(currentUser.Email, resetPasswordRequest.NewPassword, resetPasswordRequest.ConfirmPassword)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mereset kata sandi: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Reset kata sandi berhasil")
	}
}

func (h *AuthHandler) RegisterSocial() echo.HandlerFunc {
	return func(c echo.Context) error {
		var registerRequest *dto2.RegisterSocialRequest
		if err := c.Bind(&registerRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(registerRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		result, err := h.service.RegisterSocial(registerRequest)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendaftarkan akun: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Registrasi berhasil!", dto2.FormatterDetailUser(result))
	}
}

func (h *AuthHandler) LoginSocial() echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginRequest *dto2.LoginSocialRequest
		if err := c.Bind(&loginRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(loginRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		userLogin, accessToken, err := h.service.LoginSocial(loginRequest.SocialID)
		if err != nil {
			if err.Error() == "user tidak ditemukan" {
				return response.SendStatusNotFoundResponse(c, "Pengguna tidak ditemukan")
			} else if err.Error() == "akun anda belum diverifikasi" {
				return response.SendStatusUnauthorizedResponse(c, "akun anda belum diverifikasi")
			}

			logrus.Error("Kesalahan : " + err.Error())
			return response.SendStatusUnauthorizedResponse(c, "Email atau kata sandi salah")
		}

		result := &dto2.LoginResponse{
			Email:       userLogin.Email,
			AccessToken: accessToken,
		}

		return response.SendSuccessResponse(c, "Selamat datang!, Anda telah berhasil masuk.", result)
	}
}
