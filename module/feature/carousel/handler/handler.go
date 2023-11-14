package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"mime/multipart"
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

func (h *CarouselHandler) CreateCarousel() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		carouselRequest := new(dto.CreateCarouselRequest)
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
		if err := c.Bind(carouselRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(carouselRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		newCarousel := entities.CarouselModels{
			Name:  carouselRequest.Name,
			Photo: uploadedURL,
		}
		createdCarousel, err := h.service.CreateCarousel(newCarousel)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil menambahkan carousel", dto.FormatCarousel(createdCarousel))
	}
}

func (h *CarouselHandler) UpdateCarousel() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		carouselRequest := new(dto.UpdateCarouselRequest)
		carouselID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
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
		if err := c.Bind(carouselRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(carouselRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		newCarousel := entities.CarouselModels{
			ID:    carouselID,
			Name:  carouselRequest.Name,
			Photo: uploadedURL,
		}
		updatedCarousel, err := h.service.UpdateCarousel(carouselID, newCarousel)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui carousel: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil memperbarui carousel", dto.FormatCarousel(updatedCarousel))
	}
}

func (h *CarouselHandler) DeleteCarousel() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		carouselID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteCarousel(carouselID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus carousel: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil hapus carousel")
	}
}
