package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"net/http"
	"strconv"

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
		//currentUser := c.Get("CurrentUser").(*entities.UserModels)
		//if currentUser.Role != "admin" {
		//	return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		//}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 10

		var products []entities.ProductModels
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

		return response.Pagination(c, dto.FormatterProduct(products), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar produk")
	}
}

func (h *ProductHandler) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		var request dto.CreateProductRequest
		if err := c.Bind(&request); err != nil {
			c.Logger().Error("handler: invalid payload:", err.Error())
			return response.SendErrorResponse(c, http.StatusBadRequest, "Bad Request")
		}

		if err := utils.ValidateStruct(request); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		err := h.service.CreateProduct(&request)
		if err != nil {
			c.Logger().Error("handler: gagal membuat produk baru:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		return response.SendStatusCreatedResponse(c, "Product berhasil dibuat")
	}
}
