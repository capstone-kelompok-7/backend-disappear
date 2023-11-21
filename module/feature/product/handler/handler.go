package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"mime/multipart"
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
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 10

		var products []*entities.ProductModels
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

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.Pagination(c, dto.FormatterProduct(products), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Daftar produk")
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
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
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

func (h *ProductHandler) GetProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		productID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		getProductID, err := h.service.GetProductByID(productID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal mengambil produk")
		}
		return response.SendSuccessResponse(c, "Detail produk", dto.FormatProductDetail(*getProductID))
	}
}

func (h *ProductHandler) CreateProductImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan:: Anda tidak memiliki izin")
		}
		payload := new(dto.CreateProductImage)
		file, err := c.FormFile("photo")
		var uploadedURL string
		if err == nil {
			fileToUpload, err := file.Open()
			if err != nil {
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal membuka file: "+err.Error())
			}
			defer func(fileToUpload multipart.File) {
				_ = fileToUpload.Close()
			}(fileToUpload)

			uploadedURL, err = upload.ImageUploadHelper(fileToUpload)
			if err != nil {
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengunggah foto: "+err.Error())
			}
		}
		payload.Image = uploadedURL
		if err := c.Bind(payload); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang anda masukan tidak sesuai")
		}
		if err := utils.ValidateStruct(payload); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())

		}
		_, err = h.service.CreateImageProduct(*payload)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan image pada product")

	}
}

func (h *ProductHandler) GetAllProductsReview() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 10

		var products []*entities.ProductModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			products, totalItems, err = h.service.GetProductsByName(page, perPage, search)
		} else {
			products, totalItems, err = h.service.GetProductReviews(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all products:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.Pagination(c, dto.FormatReviewProductFormatter(products), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Daftar produk reviews")
	}
}

func (h *ProductHandler) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		productID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		var request dto.UpdateProduct
		if err := c.Bind(&request); err != nil {
			c.Logger().Error("handler: invalid payload:", err.Error())
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		err = h.service.UpdateProduct(productID, &request)
		if err != nil {
			c.Logger().Error("handler: gagal update produk baru:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		return response.SendStatusCreatedResponse(c, "Product berhasil diupdate")
	}
}
