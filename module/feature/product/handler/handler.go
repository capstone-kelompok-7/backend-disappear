package handler

import (
	"mime/multipart"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"

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
		pageConv := page
		perPage := 8

		var products []*entities.ProductModels
		var totalItems int64
		var err error

		search := c.QueryParam("search")
		categoryName := c.QueryParam("category_name")

		if search != "" && categoryName != "" {
			products, totalItems, err = h.service.GetProductsByCategoryAndName(categoryName, search, pageConv, perPage)
		} else if categoryName != "" {
			products, totalItems, err = h.service.GetProductsByCategoryName(categoryName, pageConv, perPage)
		} else if search != "" {
			products, totalItems, err = h.service.GetProductsByName(pageConv, perPage, search)
		} else {
			products, totalItems, err = h.service.GetAll(pageConv, perPage)
		}

		if err != nil {
			c.Logger().Error("handler: failed to fetch all products:", err.Error())
			return response.SendBadRequestResponse(c, "Gagal mendapatkan daftar produk: "+err.Error())
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterProduct(products), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar produk")
	}
}

func (h *ProductHandler) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		var request dto.CreateProductRequest
		if err := c.Bind(&request); err != nil {
			c.Logger().Error("handler: invalid payload:", err.Error())
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(request); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		createdProduct, err := h.service.CreateProduct(&request)
		if err != nil {
			c.Logger().Error("handler: gagal membuat produk baru:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan produk")
		}

		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan produk", dto.FormatCreateProductResponse(createdProduct))
	}
}

func (h *ProductHandler) GetProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		productID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}
		getProductID, err := h.service.GetProductByID(productID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail produk: "+err.Error())
		}
		getTotalSold, err := h.service.GetTotalProductSold()
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail produk: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan detail produk", dto.FormatProductDetail(*getProductID, getTotalSold))
	}
}

func (h *ProductHandler) CreateProductImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan:: Anda tidak memiliki izin")
		}
		payload := new(dto.CreateProductImage)
		file, err := c.FormFile("photo")
		var uploadedURL string
		if err == nil {
			fileToUpload, err := file.Open()
			if err != nil {
				return response.SendStatusInternalServerResponse(c, "Gagal membuka file: "+err.Error())
			}
			defer func(fileToUpload multipart.File) {
				_ = fileToUpload.Close()
			}(fileToUpload)

			uploadedURL, err = upload.ImageUploadHelper(fileToUpload)
			if err != nil {
				return response.SendStatusInternalServerResponse(c, "Gagal mengunggah foto: "+err.Error())
			}
		}
		payload.Image = uploadedURL
		if err := c.Bind(payload); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang anda masukan tidak sesuai")
		}
		if err := utils.ValidateStruct(payload); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())

		}
		image, err := h.service.CreateImageProduct(*payload)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan foto pada produk: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan foto pada produk", dto.ProductPhotoCreatedResponse(image))

	}
}

func (h *ProductHandler) GetAllProductsReview() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var products []*entities.ProductModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		filterRated := c.QueryParam("rating")

		if search != "" && filterRated != "" {
			products, totalItems, err = h.service.SearchByNameAndFilterByRating(pageConv, perPage, search, filterRated)
		} else if filterRated != "" {
			products, totalItems, err = h.service.GetRatedProductsInRange(pageConv, perPage, filterRated)
		} else if search != "" {
			products, totalItems, err = h.service.GetProductsByName(pageConv, perPage, search)
		} else {
			products, totalItems, err = h.service.GetProductReviews(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all products:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar ulasan produk")
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatReviewProductFormatter(products), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar ulasan produk")
	}
}

func (h *ProductHandler) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		productID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		var request dto.UpdateProduct
		if err := c.Bind(&request); err != nil {
			c.Logger().Error("handler: invalid payload:", err.Error())
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(request); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		_, err = h.service.UpdateProduct(productID, &request)
		if err != nil {
			c.Logger().Error("handler: gagal update produk baru:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui produk: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil memperbarui produk")
	}
}

func (h *ProductHandler) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		productId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteProduct(productId)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus produk: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil menghapus produk")
	}
}

func (h *ProductHandler) DeleteProductImageById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		productId, err := strconv.ParseUint(c.Param("idProduct"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID Product yang Anda masukkan tidak sesuai.")
		}
		imageId, err := strconv.ParseUint(c.Param("idImage"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID Product Image yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteImageProduct(productId, imageId)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus foto produk: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil menghapus foto produk")
	}
}

func (h *ProductHandler) GetAllProductsPreferences() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var products []*entities.ProductModels
		var totalItems int64
		var err error

		search := c.QueryParam("search")
		filter := c.QueryParam("filter")

		if filter != "" && search != "" {
			products, totalItems, err = h.service.GetProductsBySearchAndFilter(pageConv, perPage, filter, search)
		} else if search != "" {
			products, totalItems, err = h.service.GetProductsByName(pageConv, perPage, search)
		} else if filter != "" {
			products, totalItems, err = h.service.GetProductsByFilter(pageConv, perPage, filter)
		} else {
			products, totalItems, err = h.service.GetProductRecommendation(currentUser.ID, pageConv, perPage)
		}

		if err != nil {
			c.Logger().Error("handler: failed to fetch all products:", err.Error())
			return response.SendBadRequestResponse(c, "Gagal mendapatkan daftar produk: "+err.Error())
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterProduct(products), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar produk")
	}
}

func (h *ProductHandler) GetTopRatedProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := h.service.GetTopRatedProducts()
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail produk: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan produk lainnya", dto.FormatterOtherProduct(result))
	}
}
