package handler

import (
	"mime/multipart"
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

		var challenges []*entities.ChallengeModels
		var totalItems int64
		var err error

		search := c.QueryParam("search")
		status := c.QueryParam("status")

		if search != "" && status != "" {
			challenges, totalItems, err = h.service.GetChallengesBySearchAndStatus(page, perPage, search, status)
		} else if search != "" {
			challenges, totalItems, err = h.service.GetChallengeByTitle(page, perPage, search)
		} else if status != "" {
			challenges, totalItems, err = h.service.GetChallengeByStatus(page, perPage, status)
		} else {
			challenges, totalItems, err = h.service.GetAllChallenges(pageConv, perPage)
		}

		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar tantangan: ")
		}

		var activeChallenges []*entities.ChallengeModels
		for _, ch := range challenges {
			if ch.DeletedAt == nil {
				activeChallenges = append(activeChallenges, ch)
			}
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterChallenge(activeChallenges), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar tantangan")
	}
}

func (h *ChallengeHandler) CreateChallenge() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		challengeRequest := new(dto.CreateChallengeRequest)
		if err := c.Bind(challengeRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
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

		if err := utils.ValidateStruct(challengeRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		if challengeRequest.StartDate.After(challengeRequest.EndDate) {
			return response.SendBadRequestResponse(c, "Tanggal mulai tidak dapat setelah tanggal selesai.")
		}

		newChallenge := &entities.ChallengeModels{
			Title:       challengeRequest.Title,
			Photo:       uploadedURL,
			StartDate:   challengeRequest.StartDate,
			EndDate:     challengeRequest.EndDate,
			Description: challengeRequest.Description,
			Exp:         challengeRequest.Exp,
		}

		createdChallenge, err := h.service.CreateChallenge(newChallenge)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memambahkan tantangan: "+err.Error())
		}

		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan tantangan", dto.FormatChallenge(createdChallenge))
	}
}

func (h *ChallengeHandler) UpdateChallenge() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		challengeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}

		var updateRequest dto.UpdateChallengeRequest
		if err := c.Bind(&updateRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
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

		updatedChallenge := &entities.ChallengeModels{
			ID:          challengeID,
			Title:       updateRequest.Title,
			Photo:       uploadedURL,
			StartDate:   updateRequest.StartDate,
			EndDate:     updateRequest.EndDate,
			Description: updateRequest.Description,
			Exp:         updateRequest.Exp,
		}

		_, err = h.service.UpdateChallenge(challengeID, updatedChallenge)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui tantangan: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil memperbarui tantangan")
	}
}

func (h *ChallengeHandler) DeleteChallengeById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		challengeId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}
		err = h.service.DeleteChallenge(challengeId)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus tantangan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil menghapus tantangan")
	}
}

func (h *ChallengeHandler) GetChallengeById() echo.HandlerFunc {
	return func(c echo.Context) error {
		challengeId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}
		challenges, err := h.service.GetChallengeById(challengeId)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail tantangan: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan detail tantangan", dto.FormatChallenge(challenges))
	}
}

func (h *ChallengeHandler) CreateSubmitChallengeForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		formRequest := new(dto.CreateChallengeFormRequest)
		if err := c.Bind(formRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
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

		if err := utils.ValidateStruct(formRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		newForm := entities.ChallengeFormModels{
			ChallengeID: formRequest.ChallengeID,
			UserID:      currentUser.ID,
			Username:    formRequest.Username,
			Photo:       uploadedURL,
		}

		_, err = h.service.GetChallengeById(formRequest.ChallengeID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan informasi tantangan: "+err.Error())
		}

		createdForm, err := h.service.CreateSubmitChallengeForm(&newForm)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mengikuti tantangan: "+err.Error())
		}

		return response.SendStatusCreatedResponse(c, "Berhasil mengikuti tantangan", dto.FormatChallengeForm(createdForm))
	}
}

func (h *ChallengeHandler) GetAllSubmitChallengeForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var participants []*entities.ChallengeFormModels
		var totalItems int64
		var err error
		filterStatus := c.QueryParam("status")
		filterDate := c.QueryParam("date")

		if filterStatus != "" && filterDate != "" {
			participants, totalItems, err = h.service.GetSubmitChallengeFormByStatusAndDate(page, perPage, filterStatus, filterDate)
		} else if filterStatus != "" {
			participants, totalItems, err = h.service.GetSubmitChallengeFormByStatus(page, perPage, filterStatus)
		} else if filterDate != "" {
			participants, totalItems, err = h.service.GetSubmitChallengeFormByDateRange(page, perPage, filterDate)
		} else {
			participants, totalItems, err = h.service.GetAllSubmitChallengeForm(pageConv, perPage)
		}

		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar peserta: ")
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterChallengeForm(participants), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar peserta tantangan")
	}
}

func (h *ChallengeHandler) UpdateSubmitChallengeForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		var formRequest dto.UpdateChallengeFormStatusRequest
		formID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}
		if err := c.Bind(&formRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		validStatus := map[string]bool{"valid": true, "menunggu validasi": true, "tidak valid": true}
		if !validStatus[formRequest.Status] {
			return response.SendBadRequestResponse(c, "Status hanya bisa diubah menjadi 'valid', 'menunggu validasi', atau 'tidak valid'.")
		}

		if err := utils.ValidateStruct(&formRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		_, err = h.service.UpdateSubmitChallengeForm(formID, formRequest)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui formulir tantangan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil memperbarui formulir tantangan")
	}
}

func (h *ChallengeHandler) GetSubmitChallengeFormById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		formID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}

		form, err := h.service.GetSubmitChallengeFormById(formID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail formulir tantangan: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan detail formulir tantangan", dto.FormatChallengeForm(form))
	}
}
