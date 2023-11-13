package handler

import (
	user "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth"
	dto2 "github.com/capstone-kelompok-7/backend-disappear/module/feature/auth/dto"
	"net/http"

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
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(registerRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		existingUser, err := h.userService.GetUsersByEmail(registerRequest.Email)
		if existingUser != nil {
			return response.SendErrorResponse(c, http.StatusConflict, "Email sudah terdaftar")
		}

		newUser := &user.UserModels{
			Email:    registerRequest.Email,
			Password: registerRequest.Password,
		}

		_, err = h.service.Register(newUser)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Registrasi berhasil! Silakan cek email anda untuk OTP.")
	}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginRequest dto2.LoginRequest
		if err := c.Bind(&loginRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(loginRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		userLogin, accessToken, err := h.service.Login(loginRequest.Email, loginRequest.Password)
		if err != nil {
			if err.Error() == "user tidak ditemukan" {
				return response.SendErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
			} else if err.Error() == "akun anda belum diverifikasi" {
				return response.SendErrorResponse(c, http.StatusNotFound, "akun anda belum diverifikasi")
			}

			logrus.Error("Kesalahan : " + err.Error())
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Email atau kata sandi salah")
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
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		if err := utils.ValidateStruct(emailRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		if err := h.service.VerifyEmail(emailRequest.Email, emailRequest.OTP); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Email berhasil diverifikasi!, Silahkan login.")
	}
}

func (h *AuthHandler) ResendOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		var emailRequest dto2.ResendOTPRequest

		if err := c.Bind(&emailRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(emailRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		newOTP, err := h.service.ResendOTP(emailRequest.Email)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal mengirim OTP: "+err.Error())
		}

		err = email.EmaiilService(emailRequest.Email, newOTP.OTP)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal mengirim OTP ke email: "+err.Error())

		}

		return response.SendStatusOkResponse(c, "OTP berhasil dikirim kembali!, Silahkan cek email anda.")

	}
}

func (h *AuthHandler) ForgotPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var forgotPasswordRequest dto2.ForgotPasswordRequest
		if err := c.Bind(&forgotPasswordRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(forgotPasswordRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		newOTP, err := h.service.ResendOTP(forgotPasswordRequest.Email)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal mengirim OTP: "+err.Error())
		}

		err = email.EmaiilService(forgotPasswordRequest.Email, newOTP.OTP)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal mengirim OTP ke email: "+err.Error())

		}

		return response.SendStatusOkResponse(c, "Permintaan OTP berhasil dikirim")
	}
}

func (h *AuthHandler) VerifyOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		var emailRequest dto2.EmailRequest

		if err := c.Bind(&emailRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		if err := utils.ValidateStruct(emailRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		accessToken, err := h.service.VerifyOTP(emailRequest.Email, emailRequest.OTP)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal verifikasi OTP: "+err.Error())
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
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal Mengikat Data: Pengikatan data ke struktur gagal")
		}

		if err := utils.ValidateStruct(resetPasswordRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		currentUser := c.Get("CurrentUser").(*user.UserModels)
		err := h.service.ResetPassword(currentUser.Email, resetPasswordRequest.NewPassword, resetPasswordRequest.ConfirmPassword)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mereset kata sandi: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Reset kata sandi berhasil")
	}
}
