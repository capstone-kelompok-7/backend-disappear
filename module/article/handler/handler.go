package handler

import (
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/article"
	"github.com/capstone-kelompok-7/backend-disappear/module/article/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	service article.ServiceArticleInterface
}

func NewArticleHandler(service article.ServiceArticleInterface) article.HandlerArticleInterface {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) GetAllArticles() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 10

		var articles []domain.Articles
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			articles, totalItems, err = h.service.GetArticlesByTitle(page, perPage, search)
		} else {
			articles, totalItems, err = h.service.GetAll(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all articles:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, domain.FormatterArticle(articles), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar artikel")
	}
}
