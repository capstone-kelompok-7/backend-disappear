package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"mime/multipart"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	service category.ServiceCategoryInterface
}

func NewCategoryHandler(service category.ServiceCategoryInterface) category.HandlerCategoryInterface {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) CreateCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		categoryRequest := new(dto.CreateCategoryRequest)
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
		if err := c.Bind(categoryRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(categoryRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		newCategory := &entities.CategoryModels{
			Name:  categoryRequest.Name,
			Photo: uploadedURL,
		}
		createdCategory, err := h.service.CreateCategory(newCategory)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan kategori: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan kategori", dto.FormatCategory(createdCategory))
	}
}

func (h *CategoryHandler) GetAllCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var categories []*entities.CategoryModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			categories, totalItems, err = h.service.GetCategoryByName(page, perPage, search)
		} else {
			categories, totalItems, err = h.service.GetAllCategory(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all carousels:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar kategori: ")
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterCategory(categories), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar kategori")
	}
}

func (h *CategoryHandler) GetCategoryById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		id := c.Param("id")
		categoryID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}

		getCategory, err := h.service.GetCategoryById(categoryID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail kategori: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Data kategori", getCategory)
	}
}

func (h *CategoryHandler) UpdateCategoryById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		updateRequest := new(dto.UpdateCategoryRequest)
		categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}

		existingCategory, err := h.service.GetCategoryById(categoryID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan kategori: "+err.Error())
		}
		if existingCategory == nil {
			return response.SendStatusNotFoundResponse(c, "Kategori tidak ditemukan: "+err.Error())
		}
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
		if err := c.Bind(updateRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(updateRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		newData := &entities.CategoryModels{
			ID:    categoryID,
			Name:  updateRequest.Name,
			Photo: uploadedURL,
		}
		err = h.service.UpdateCategoryById(categoryID, newData)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui kategori: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil mengubah kategori")
	}
}

func (h *CategoryHandler) DeleteCategoryById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}

		err = h.service.DeleteCategoryById(categoryID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus kategori: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil menghapus kategori")
	}
}
