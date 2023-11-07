package handler

import (
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/auth/domain"
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	user "github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
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

		registerRequest := new(domain.RegisterRequest)
		if err := c.Bind(registerRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(registerRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		_, err := h.userService.GetUsersByEmail(registerRequest.Email)
		fmt.Println(registerRequest.Email)
		if err == nil {
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
		return response.SendStatusOkResponse(c, "Registrasi berhasil! Silakan masuk untuk memulai.")
	}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginRequest domain.LoginRequest
		if err := c.Bind(&loginRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(loginRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		userLogin, accessToken, err := h.service.Login(loginRequest.Email, loginRequest.Password)
		if err != nil {
			if err.Error() == "user not found" {
				return response.SendErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
			}
			logrus.Error("Kesalahan : " + err.Error())
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Email atau kata sandi salah")
		}

		result := &domain.LoginResponse{
			Email:       userLogin.Email,
			AccessToken: accessToken,
		}

		return response.SendSuccessResponse(c, "Selamat datang!, Anda telah berhasil masuk.", result)
	}
}
