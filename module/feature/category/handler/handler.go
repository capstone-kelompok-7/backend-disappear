package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/category/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"mime/multipart"
	"net/http"
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
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan:: Anda tidak memiliki izin")
		}
		categoryRequest := new(dto.CreateCategoryRequest)
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
		if err := c.Bind(categoryRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(categoryRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		newCategory := &entities.CategoryModels{
			Name:  categoryRequest.Name,
			Photo: uploadedURL,
		}
		createdCategory, err := h.service.CreateCategory(newCategory)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil menambahkan kategory", dto.FormatCategory(createdCategory))
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
			if err != nil {
				c.Logger().Error("handler: failed to fetch categories by name:", err.Error())
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
			}
		} else {
			categories, totalItems, err = h.service.GetAllCategory(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all categories:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, dto.FormatterCategory(categories), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar kategori")
	}
}

func (h *CategoryHandler) GetCategoryByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, perPage := 1, 10
		name := c.Param("name")

		categories, totalItems, err := h.service.GetCategoryByName(page, perPage, name)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil kategori")
		}

		if len(categories) == 0 {
			return response.SendErrorResponse(c, http.StatusNotFound, "Kategori tidak ditemukan")
		}

		return response.Pagination(c, dto.FormatterCategory(categories), 1, 1, int(totalItems), 1, 1, "Daftar kategori")
	}
}

func (h *CategoryHandler) UpdateCategoryById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan:: Anda tidak memiliki izin")
		}
		updateRequest := new(dto.UpdateCategoryRequest)
		categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		existingCategory, err := h.service.GetCategoryById(categoryID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan kategori: "+err.Error())
		}
		if existingCategory == nil {
			return response.SendErrorResponse(c, http.StatusNotFound, "Kategori tidak ditemukan"+err.Error())
		}
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
		if err := c.Bind(updateRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(updateRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		newData := &entities.CategoryModels{
			ID:    categoryID,
			Name:  updateRequest.Name,
			Photo: uploadedURL,
		}
		updatedCategory, err := h.service.UpdateCategoryById(categoryID, newData)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal update category:"+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mengubah kategori", dto.FormatCategory(updatedCategory))
	}
}

func (h *CategoryHandler) DeleteCategoryById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan:: Anda tidak memiliki izin")
		}

		categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		existingCategory, err := h.service.GetCategoryById(categoryID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan kategori: "+err.Error())
		}
		if existingCategory == nil {
			return response.SendErrorResponse(c, http.StatusNotFound, "Kategori tidak ditemukan"+err.Error())
		}

		err = h.service.DeleteCategoryById(categoryID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Server Internal Error"+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil hapus kategori")
	}
}
