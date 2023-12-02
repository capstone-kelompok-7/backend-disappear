package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SuccessResponses struct {
	Message      string      `json:"message"`
	ExtraMessage string      `json:"label"`
	Data         interface{} `json:"data"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalItems  int `json:"total_items"`
	NextPage    int `json:"next_page"`
	PrevPage    int `json:"prev_page"`
}

type PaginationRes struct {
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

func SendStatusForbiddenResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusForbidden, ErrorResponse{
		Message: message,
	})
}

func SendStatusCreatedResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func SendStatusOkResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, ErrorResponse{
		Message: message,
	})
}

func SendStatusOkWithDataResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func SendStatusOkWithDataResponses(c echo.Context, message string, extraMessage string, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessResponses{
		Message:      message,
		ExtraMessage: extraMessage,
		Data:         data,
	})
}

func SendPaginationResponse(c echo.Context, data interface{}, currentPage, totalPages, totalItems, nextPage, prevPage int, message string) error {
	pagination := PaginationMeta{
		CurrentPage: currentPage,
		TotalPage:   totalPages,
		TotalItems:  totalItems,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
	return c.JSON(http.StatusOK, PaginationRes{
		Message: message,
		Data:    data,
		Meta:    pagination,
	})
}

func SendStatusInternalServerResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Message: message,
	})
}

func SendBadRequestResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Message: message,
	})
}

func SendSuccessResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func SendStatusConflictResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusConflict, ErrorResponse{
		Message: message,
	})
}

func SendStatusNotFoundResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, ErrorResponse{
		Message: message,
	})
}

func SendStatusUnauthorizedResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, ErrorResponse{
		Message: message,
	})
}
