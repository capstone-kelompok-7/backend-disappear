package handler

import (
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
			return response.SendErrorResponse(c, http.StatusBadRequest, "Error Binding Data")
		}

		if err := utils.ValidateStruct(registerRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validation failed : "+err.Error())
		}

		_, err := h.userService.GetUsersByEmail(registerRequest.Email)
		if err == nil {
			return response.SendErrorResponse(c, http.StatusConflict, "Email already registered")
		}

		newUser := &user.UserModels{
			Email:    registerRequest.Email,
			Phone:    registerRequest.Phone,
			Password: registerRequest.Password,
		}

		_, err = h.service.Register(newUser)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error"+err.Error())
		}
		return response.SendStatusOkResponse(c, "Success Register, Please verify your account via email")
	}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginRequest domain.LoginRequest
		if err := c.Bind(&loginRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Error Binding Data")
		}

		if err := utils.ValidateStruct(loginRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
		}

		userLogin, accessToken, err := h.service.Login(loginRequest.Email, loginRequest.Password)
		if err != nil {
			logrus.Error("Error : " + err.Error())
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Invalid email or password ")
		}

		result := &domain.LoginResponse{
			Email:       userLogin.Email,
			AccessToken: accessToken,
		}

		return response.SendSuccessResponse(c, "Success", result)
	}
}
