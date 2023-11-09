package handler

import (
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/category"
	"github.com/capstone-kelompok-7/backend-disappear/module/category/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService category.ServiceCategoryInterface
}

func NewCategoryHandler(categoryService category.ServiceCategoryInterface) category.HandlerCategoryInterface {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		var categoryData domain.CategoryModels
		if err := c.Bind(&categoryData); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Gagal memparsing data kategori"})
		}

		createdCategory, err := h.categoryService.CreateCategory(&categoryData)

		if err != nil {
			if err.Error() == "Kategori dengan nama yang sama sudah ada" {
				return c.JSON(http.StatusConflict, map[string]interface{}{"message": "Kategori sudah ada"})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Gagal membuat kategori"})
			}
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{"message": "Kategori berhasil dibuat", "data": createdCategory})
	}
}

func (h *CategoryHandler) GetAllCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 10

		var categories []*domain.CategoryModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			categories, totalItems, err = h.categoryService.GetCategoryByName(page, perPage, search)
			if err != nil {
				c.Logger().Error("handler: failed to fetch categories by name:", err.Error())
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
			}
		} else {
			categories, totalItems, err = h.categoryService.GetAllCategory(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all categories:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		current_page, total_pages := h.categoryService.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.categoryService.GetNextPage(current_page, total_pages)
		prevPage := h.categoryService.GetPrevPage(current_page)

		return response.Pagination(c, domain.FormatterCategory(categories), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar kategori")
	}
}

func (h *CategoryHandler) GetCategoryByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, perPage := 1, 10
		name := c.Param("name")

		categories, totalItems, err := h.categoryService.GetCategoryByName(page, perPage, name)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil kategori")
		}

		if len(categories) == 0 {
			return response.SendErrorResponse(c, http.StatusNotFound, "Kategori tidak ditemukan")
		}

		return response.Pagination(c, domain.FormatterCategory(categories), 1, 1, int(totalItems), 1, 1, "Daftar kategori")
	}
}

func (h *CategoryHandler) UpdateCategoryById() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		}

		var updatedCategoryData domain.CategoryModels
		if err := c.Bind(&updatedCategoryData); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Gagal memparsing data kategori")
		}

		updatedCategory, err := h.categoryService.UpdateCategoryById(id, &updatedCategoryData)

		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui kategori")
		}

		if updatedCategory == nil {
			return response.SendErrorResponse(c, http.StatusNotFound, "Kategori tidak ditemukan")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Kategori berhasil diperbarui", "data": updatedCategory})
	}
}

func (h *CategoryHandler) DeleteCategoryById() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		}

		err = h.categoryService.DeleteCategoryById(id)

		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus kategori")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Kategori berhasil dihapus"})
	}
}
