package handler

import (
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

		var challenge []*entities.ChallengeModels
		var err error
		search := c.QueryParam("search")
		if search != "" {
			challenge, _, err = h.service.GetChallengeByTitle(page, perPage, search)
		} else {
			challenge, _, err = h.service.GetAllChallenges(pageConv, perPage)
		}
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		var activeChallenges []*entities.ChallengeModels
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
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}

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

		updatedChallenge := &entities.ChallengeModels{
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

func (h *ChallengeHandler) CreateSubmitChallengeForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		formRequest := new(dto.CreateChallengeFormRequest)
		if err := c.Bind(formRequest); err != nil {
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

		if err := utils.ValidateStruct(formRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		newForm := entities.ChallengeFormModels{
			ChallengeID: formRequest.ChallengeID,
			UserID:      currentUser.ID,
			Username:    formRequest.Username,
			Photo:       uploadedURL,
		}

		challenge, err := h.service.GetChallengeById(formRequest.ChallengeID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan informasi tantangan: "+err.Error())
		}

		createdForm, err := h.service.CreateSubmitChallengeForm(&newForm)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengikuti tantangan: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mengikuti tantangan", dto.FormatChallengeForm(createdForm, challenge.Exp, createdForm.CreatedAt))
	}
}

func (h *ChallengeHandler) GetAllSubmitChallengeForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var partisipan []*entities.ChallengeFormModels
		var totalItems int64
		var err error
		filterStatus := c.QueryParam("status")

		validStatus := map[string]bool{"valid": true, "menunggu validasi": true, "tidak valid": true}
		if filterStatus != "" && !validStatus[filterStatus] {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Invalid status parameter")
		}

		if filterStatus != "" {
			partisipan, totalItems, err = h.service.GetSubmitChallengeFormByStatus(pageConv, perPage, filterStatus)
		} else {
			partisipan, totalItems, err = h.service.GetAllSubmitChallengeForm(pageConv, perPage)
		}

		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
		}

		challenge, _ := h.service.GetChallengeById(partisipan[0].ChallengeID)
		exp := challenge.Exp
		createdAt := challenge.CreatedAt

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, dto.FormatterChallengeForm(partisipan, exp, createdAt), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar peserta")
	}
}

func (h *ChallengeHandler) UpdateSubmitChallengeForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		var formRequest dto.UpdateChallengeFormStatusRequest
		formID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		if err := c.Bind(&formRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		validStatus := map[string]bool{"valid": true, "menunggu validasi": true, "tidak valid": true}
		if !validStatus[formRequest.Status] {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Status hanya bisa diubah menjadi 'valid', 'menunggu validasi', atau 'tidak valid'.")
		}

		if err := utils.ValidateStruct(&formRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		_, err = h.service.UpdateSubmitChallengeForm(formID, formRequest)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "form berhasil dirubah")
	}
}

func (h *ChallengeHandler) GetSubmitChallengeFormById() echo.HandlerFunc {
	return func(c echo.Context) error {
		formID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		form, err := h.service.GetSubmitChallengeFormById(formID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan form: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan form submit", dto.FormatChallengeForm(form, form.Exp, form.CreatedAt))
	}
}
