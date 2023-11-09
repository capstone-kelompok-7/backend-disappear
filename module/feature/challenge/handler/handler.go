package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/dto"
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type ChallengeHandler struct {
	service challenge.ServiceChallengeInterface
}

func NewChallengeHandler(service challenge.ServiceChallengeInterface) challenge.HandlerChallengeInterface {
	return &ChallengeHandler{
		service: service,
	}
}

func (h *ChallengeHandler) GetAllChallenges() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var challenge []entities.ChallengeModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			challenge, totalItems, err = h.service.GetChallengeByTitle(page, perPage, search)
		} else {
			challenge, totalItems, err = h.service.GetAllChallenges(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all challenge:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, dto.FormatterChallenge(challenge), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar challenge")
	}
}
