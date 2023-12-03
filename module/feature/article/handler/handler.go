package handler

import (
	"mime/multipart"
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
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		articleRequest := new(dto.CreateArticleRequest)
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

		if err := c.Bind(articleRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai: "+err.Error())
		}

		if err := utils.ValidateStruct(articleRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		newArticle := &entities.ArticleModels{
			Title:   articleRequest.Title,
			Photo:   uploadedURL,
			Content: articleRequest.Content,
		}

		createdArticle, err := h.service.CreateArticle(newArticle)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan artikel: "+err.Error())
		}

		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan artikel", dto.FormatArticle(*createdArticle))
	}
}

func (h *ArticleHandler) UpdateArticleById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		updateRequest := new(dto.UpdateArticleRequest)
		articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai: "+err.Error())
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
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai: "+err.Error())
		}

		if err := utils.ValidateStruct(updateRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		newData := &entities.ArticleModels{
			Title:   updateRequest.Title,
			Photo:   uploadedURL,
			Content: updateRequest.Content,
		}

		_, err = h.service.UpdateArticleById(articleID, newData)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui artikel: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil memperbarui artikel")
	}
}

func (h *ArticleHandler) DeleteArticleById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai: "+err.Error())
		}

		err = h.service.DeleteArticleById(articleID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus artikel: "+err.Error())
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
				return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan pencarian artikel: "+err.Error())
			}
		} else if filterType != "" {
			articles, totalItems, err = h.service.GetArticlesByDateRange(page, perPage, filterType)
			if err != nil {
				return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan filter artikel: "+err.Error())
			}
		} else {
			articles, totalItems, err = h.service.GetAll(pageConv, perPage)
			if err != nil {
				return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan artikel: "+err.Error())
			}
		}

		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar artikel: "+err.Error())
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterArticle(articles), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar artikel")
	}
}

func (h *ArticleHandler) GetArticleById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		articleID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai: "+err.Error())
		}

		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		incrementViews := currentUser.Role != "admin"

		getArticleID, err := h.service.GetArticleById(articleID, incrementViews)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail artikel: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan detail artikel", dto.FormatArticle(*getArticleID))
	}
}

func (h *ArticleHandler) BookmarkArticle() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		bookmark := new(dto.UserBookmarkRequest)
		if err := c.Bind(bookmark); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang anda masukkan tidak sesuai.")
		}
		if err := utils.ValidateStruct(bookmark); err != nil {
			return response.SendStatusInternalServerResponse(c, "Validasi gagal: "+err.Error())
		}
		newBookmark := &entities.UserBookmarkModels{
			UserID: currentUser.ID,
			ArticleID: bookmark.ArticleID,
		}
		if err := h.service.BookmarkArticle(newBookmark); err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menyimpan artikel: " + err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil menyimpan artikel")
	}
}

func (h *ArticleHandler) GetUsersBookmark() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		result, err := h.service.GetUserBookmarkArticle(currentUser.ID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan data artikel tersimpan user : "+err.Error())
		}

		formattedResponse, err := dto.UserBookmarkFormatter(result)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memformat response : "+ err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan data artikel tersimpan user", formattedResponse)
	}
}