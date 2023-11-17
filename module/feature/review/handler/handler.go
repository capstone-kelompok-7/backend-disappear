package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/review/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"net/http"
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
		reviewRequest := new(dto.CreateReviewRequest)
		if err := c.Bind(reviewRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(reviewRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		newReview := &entities.ReviewModels{
			UserID:      currentUser.ID,
			ProductID:   reviewRequest.ProductID,
			Rating:      reviewRequest.Rating,
			Description: reviewRequest.Description,
		}
		createdAddress, err := h.service.CreateReview(newReview)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil menambahkan review", dto.FormatReview(createdAddress))
	}
}
