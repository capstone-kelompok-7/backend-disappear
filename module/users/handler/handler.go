package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	service users.ServiceUserInterface
}

func NewUserHandler(service users.ServiceUserInterface) users.HandlerUserInterface {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetAllUsers() echo.HandlerFunc {
	return func(c echo.Context) error {

		result, err := h.service.GetAllUsers()
		if err != nil {
			c.Logger().Error("handler: failed to fetch all users:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}
		if len(result) == 0 {
			return response.SendSuccessResponse(c, "Success", nil)
		} else {
			return response.SendSuccessResponse(c, "Success", result)
		}
	}
}

func (h *UserHandler) GetUsersByEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.QueryParam("email")
		if email == "" {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Email parameter is missing")
		}
		user, err := h.service.GetUsersByEmail(email)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusNotFound, "User not found")
		}

		return response.SendSuccessResponse(c, "Success", user)
	}
}

func (h *UserHandler) ChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req domain.ChangePasswordRequest
		if err := c.Bind(&req); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal Mengikat Data: Pengikatan data ke struktur gagal")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		user, err := h.service.ChangePassword(req.Email, req.OldPassword, req.NewPassword)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengganti kata sandi: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Kata sandi berhasil diubah", user)
	}
}
