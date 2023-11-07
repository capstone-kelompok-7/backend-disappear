package handler

import (
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/product/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service product.ServiceProductInterface
}

func NewProductHandler(service product.ServiceProductInterface) product.HandlerProductInterface {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 10

		var products []domain.ProductModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			products, totalItems, err = h.service.GetProductsByName(page, perPage, search)
		} else {
			products, totalItems, err = h.service.GetAll(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all products:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, domain.FormatterProduct(products), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar produk")
	}
}
