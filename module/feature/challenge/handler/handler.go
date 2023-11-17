package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"

	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
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
		var err error
		search := c.QueryParam("search")
		if search != "" {
			challenge, _, err = h.service.GetChallengeByTitle(page, perPage, search)
		} else {
			challenge, _, err = h.service.GetAllChallenges(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all challenge:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		var activeChallenges []entities.ChallengeModels
		for _, ch := range challenge {
			if ch.DeletedAt == nil {
				activeChallenges = append(activeChallenges, ch)
			}
		}

		totalItems := int64(len(activeChallenges))

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, len(activeChallenges), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, dto.FormatterChallenge(activeChallenges), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar challenge")
	}
}

func (h *ChallengeHandler) CreateChallenge() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		challengeRequest := new(dto.CreateChallengeRequest)
		if err := c.Bind(challengeRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		file, err := c.FormFile("photo")
		var uploadedURL string
		if err == nil {
			fileToUpload, err := file.Open()
			if err != nil {
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal membuka file: "+err.Error())
			}
			defer fileToUpload.Close()

			uploadedURL, err = upload.ImageUploadHelper(fileToUpload)
			if err != nil {
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengunggah foto: "+err.Error())
			}
		}

		if err := utils.ValidateStruct(challengeRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		if challengeRequest.StartDate.After(challengeRequest.EndDate) {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Tanggal mulai tidak dapat setelah tanggal selesai.")
		}

		newChallenge := entities.ChallengeModels{
			Title:       challengeRequest.Title,
			Photo:       uploadedURL,
			StartDate:   challengeRequest.StartDate,
			EndDate:     challengeRequest.EndDate,
			Description: challengeRequest.Description,
			Exp:         challengeRequest.Exp,
		}

		// Log nilai waktu sebelum disimpan
		log.Printf("Waktu sebelum disimpan: StartDate %v, EndDate %v", newChallenge.StartDate, newChallenge.EndDate)

		createdChallenge, err := h.service.CreateChallenge(newChallenge)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}

		// Log nilai waktu setelah disimpan
		log.Printf("Waktu setelah disimpan: StartDate %v, EndDate %v", createdChallenge.StartDate, createdChallenge.EndDate)

		return response.SendSuccessResponse(c, "Berhasil menambahkan tantangan", dto.FormatChallenge(createdChallenge))
	}
}

func (h *ChallengeHandler) UpdateChallenge() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		challengeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "ID tantangan tidak valid")
		}

		var updateRequest dto.UpdateChallengeRequest
		if err := c.Bind(&updateRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		file, err := c.FormFile("photo")
		var uploadedURL string
		if err == nil {
			fileToUpload, err := file.Open()
			if err != nil {
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal membuka file: "+err.Error())
			}
			defer fileToUpload.Close()

			uploadedURL, err = upload.ImageUploadHelper(fileToUpload)
			if err != nil {
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengunggah foto: "+err.Error())
			}
		}

		updatedChallenge := entities.ChallengeModels{
			ID:          challengeID,
			Title:       updateRequest.Title,
			Photo:       uploadedURL,
			StartDate:   updateRequest.StartDate,
			EndDate:     updateRequest.EndDate,
			Description: updateRequest.Description,
			Exp:         updateRequest.Exp,
		}

		updatedChallenge, err = h.service.UpdateChallenge(challengeID, updatedChallenge)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil memperbarui tantangan", dto.FormatChallenge(updatedChallenge))
	}
}

func (h *ChallengeHandler) DeleteChallengeById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		challengeId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteChallenge(challengeId)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus tantangan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil hapus tantangan")
	}
}

func (h *ChallengeHandler) GetChallengeById() echo.HandlerFunc {
	return func(c echo.Context) error {
		challengeId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		challenges, err := h.service.GetChallengeById(challengeId)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan tantangan: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan tantangan", dto.FormatChallenge(challenges))
	}
}
