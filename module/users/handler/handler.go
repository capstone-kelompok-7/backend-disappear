package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	user "github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

func (h *UserHandler) GetUsersById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*user.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized: You don't have permission")
		}
		id := c.Param("id")
		if id == "" {
			return response.SendErrorResponse(c, http.StatusBadRequest, "ID parameter is missing")
		}

		userID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		}

		user, err := h.service.GetUsersById(userID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusNotFound, "User not found")
		}

		return response.SendSuccessResponse(c, "Success", user)
	}
}
