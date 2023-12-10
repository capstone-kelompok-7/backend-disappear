package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/carousel/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"mime/multipart"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type CarouselHandler struct {
	service carousel.ServiceCarouselInterface
}

func NewCarouselHandler(service carousel.ServiceCarouselInterface) carousel.HandlerCarouselInterface {
	return &CarouselHandler{
		service: service,
	}
}

func (h *CarouselHandler) GetAllCarousels() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var carousels []*entities.CarouselModels
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
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar carausel: "+err.Error())
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterCarousel(carousels), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar carousel")
	}
}

func (h *CarouselHandler) GetCarouselById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		id := c.Param("id")
		carouselId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}
		getCarousel, err := h.service.GetCarouselById(carouselId)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail carousel: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Data carousel", getCarousel)
	}
}

func (h *CarouselHandler) CreateCarousel() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.CreateCarouselRequest)
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
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		newCarousel := &entities.CarouselModels{
			Name:  req.Name,
			Photo: uploadedURL,
		}
		createdCarousel, err := h.service.CreateCarousel(newCarousel)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan carousel: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan carousel", dto.FormatCarousel(createdCarousel))
	}
}

func (h *CarouselHandler) UpdateCarousel() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.UpdateCarouselRequest)
		carouselID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
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
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		newCarousel := &entities.CarouselModels{
			ID:    carouselID,
			Name:  req.Name,
			Photo: uploadedURL,
		}
		err = h.service.UpdateCarousel(carouselID, newCarousel)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui carousel: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil memperbarui carousel")
	}
}

func (h *CarouselHandler) DeleteCarousel() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		carouselID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}
		err = h.service.DeleteCarousel(carouselID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus carousel: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil hapus carousel")
	}
}
