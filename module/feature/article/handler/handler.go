package handler

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/article"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/article/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"

	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	service article.ServiceArticleInterface
}

func NewArticleHandler(service article.ServiceArticleInterface) article.HandlerArticleInterface {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) CreateArticle() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		articleRequest := new(dto.CreateArticleRequest)
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

		if err := c.Bind(articleRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai: "+err.Error())
		}

		if err := utils.ValidateStruct(articleRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		newArticle := &entities.ArticleModels{
			Title:   articleRequest.Title,
			Photo:   uploadedURL,
			Content: articleRequest.Content,
		}

		createdArticle, err := h.service.CreateArticle(newArticle)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil menambahkan artikel", dto.FormatArticle(*createdArticle))
	}
}

func (h *ArticleHandler) UpdateArticleById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		updateRequest := new(dto.UpdateArticleRequest)
		articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai: "+err.Error())
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
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai: "+err.Error())
		}

		if err := utils.ValidateStruct(updateRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		newData := &entities.ArticleModels{
			Title:   updateRequest.Title,
			Photo:   uploadedURL,
			Content: updateRequest.Content,
		}

		updatedArticle, err := h.service.UpdateArticleById(articleID, newData)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengubah artikel: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mengubah artikel", dto.FormatArticle(*updatedArticle))
	}
}

func (h *ArticleHandler) DeleteArticleById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai: "+err.Error())
		}

		err = h.service.DeleteArticleById(articleID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus artikel: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil menghapus artikel")
	}
}

func (h *ArticleHandler) GetAllArticles() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 10

		var articles []entities.ArticleModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		filterType := c.QueryParam("filter_type")
		if search != "" {
			articles, totalItems, err = h.service.GetArticlesByTitle(page, perPage, search)
			if err != nil {
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
			}
		} else if filterType != "" {
			articles, totalItems, err = h.service.GetArticlesByDateRange(page, perPage, filterType)
        	if err != nil {
            	return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
        	}
		} else {
			articles, totalItems, err = h.service.GetAll(pageConv, perPage)
			if err != nil {
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
			}
		}

		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
		}

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, dto.FormatterArticle(articles), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar artikel")
	}
}

func (h *ArticleHandler) GetArticleById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		articleID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai: " + err.Error())
		}

		currentUser := c.Get("CurrentUser").(*entities.UserModels)
        incrementViews := currentUser.Role != "admin"

		getArticleID, err := h.service.GetArticleById(articleID, incrementViews)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil artikel: " + err.Error())
		}

		return response.SendSuccessResponse(c, "Detail artikel", dto.FormatArticle(*getArticleID))
	}
}