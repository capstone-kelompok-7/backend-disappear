package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel/dto"
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type CarouselHandler struct {
	service carousel.ServiceCarouselInterface
}

func NewCarouselHandler(service carousel.ServiceCarouselInterface) carousel.HandlerCarouselInterface {
	return &CarouselHandler{service: service}
}

func (h *CarouselHandler) GetAllCarousels() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var carousels []entities.CarouselModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			carousels, totalItems, err = h.service.GetCarouselsByName(page, perPage, search)
		} else {
			carousels, totalItems, err = h.service.GetAll(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all carousels:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, dto.FormatterCarousel(carousels), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar carousel")
	}
}
