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

func SendErrorResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, ErrorResponse{
		Message: message,
	})
}

func SendSuccessResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}
func SendStatusOkResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, ErrorResponse{
		Message: message,
	})
}

func PaginationResponse(c echo.Context, data interface{}, totalItems, page, pageSize int, message string) error {
	pagination := map[string]interface{}{
		"totalItems": totalItems,
		"page":       page,
		"pageSize":   pageSize,
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    map[string]interface{}{"items": data, "pagination": pagination},
	})
}

func Pagination(c echo.Context, data interface{}, currentPage, totalPages, totalItems, nextPage, prevPage int, message string) error {
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
