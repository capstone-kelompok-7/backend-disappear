package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"strconv"
)

type ReviewHandler struct {
	service review.ServiceReviewInterface
}

func NewReviewHandler(service review.ServiceReviewInterface) review.HandlerReviewInterface {
	return &ReviewHandler{
		service: service,
	}
}

func (h *ReviewHandler) CreateReview() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		reviewRequest := new(dto.CreateReviewRequest)
		if err := c.Bind(reviewRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(reviewRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		newReview := &entities.ReviewModels{
			UserID:      currentUser.ID,
			ProductID:   reviewRequest.ProductID,
			Rating:      reviewRequest.Rating,
			Description: reviewRequest.Description,
		}
		result, err := h.service.CreateReview(newReview)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan ulasan "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan ulasan", dto.CreateFormatReview(result))
	}
}

func (h *ReviewHandler) CreateReviewImages() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		reviewRequest := new(dto.CreatePhotoReviewRequest)
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
		if err := c.Bind(reviewRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(reviewRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		newReviewPhoto := &entities.ReviewPhotoModels{
			ReviewID: reviewRequest.ReviewID,
			ImageURL: uploadedURL,
		}
		result, err := h.service.CreateReviewImages(newReviewPhoto)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan foto ulasan "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan foto ulasan", dto.FormatReviewPhoto(result))
	}
}

func (h *ReviewHandler) GetReviewById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		reviewID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}
		result, err := h.service.GetReviewById(reviewID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapat ulasan: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan detail ulasan", dto.FormatReview(result))
	}
}

func (h *ReviewHandler) GetDetailReviewProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		productID := c.Param("id")
		id, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}

		result, err := h.service.GetReviewsProductByID(id)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan ulasan produk: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan ulasan produk", dto.FormatProductDetail(result))
	}
}
