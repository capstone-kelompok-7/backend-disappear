package handler

import (
	"mime/multipart"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service users.ServiceUserInterface
}

func NewUserHandler(service users.ServiceUserInterface) users.HandlerUserInterface {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetUsersByEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		email := c.QueryParam("email")
		if email == "" {
			return response.SendBadRequestResponse(c, "Format email yang Anda masukkan tidak sesuai.")
		}
		user, err := h.service.GetUsersByEmail(email)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan data pengguna: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapat data pengguna", user)
	}
}

func (h *UserHandler) ChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var updateRequest dto.UpdatePasswordRequest
		if err := c.Bind(&updateRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format email yang Anda masukkan tidak sesuai.")
		}
		if err := utils.ValidateStruct(updateRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		currentUser := c.Get("CurrentUser").(*entities.UserModels)

		err := h.service.ValidatePassword(currentUser.ID, updateRequest.OldPassword, updateRequest.NewPassword, updateRequest.ConfirmPassword)
		if err != nil {
			return response.SendStatusConflictResponse(c, "Kata sandi lama salah: "+err.Error())
		}
		err = h.service.ChangePassword(currentUser.ID, updateRequest)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui kata sandi:: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil memperbarui kata sandi")
	}
}

func (h *UserHandler) GetUsersById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}

		userID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		user, err := h.service.GetUsersById(userID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail pengguna: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapat detail pengguna", dto.FormatterDetailUser(user))
	}
}

func (h *UserHandler) GetAllUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var user []*entities.UserModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		filter := c.QueryParam("filter")

		switch {
		case search != "" && filter != "":
			user, totalItems, err = h.service.GetUsersBySearchAndFilter(page, perPage, search, filter)
		case search != "":
			user, totalItems, err = h.service.GetUsersByName(page, perPage, search)
		case filter != "":
			user, totalItems, err = h.service.GetUsersByLevel(page, perPage, filter)
		default:
			user, totalItems, err = h.service.GetAllUsers(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all users:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar customer: "+err.Error())
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterUsersPagination(user), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar customer")
	}
}

func (h *UserHandler) EditProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		var editProfileRequest dto.EditProfileRequest
		if err := c.Bind(&editProfileRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}
		if err := utils.ValidateStruct(editProfileRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
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
			editProfileRequest.PhotoProfile = uploadedURL
		}
		_, err = h.service.EditProfile(currentUser.ID, editProfileRequest)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui profil: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Profil berhasil diperbarui")
	}
}

func (h *UserHandler) DeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteAccount(userID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus pengguna: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil hapus pengguna")
	}
}

func (h *UserHandler) GetLeaderboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		limit := 5
		result, err := h.service.GetLeaderboardByExp(limit)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan leaderboard: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan leaderboard", dto.FormatterUserLeaderboard(result))
	}
}

func (h *UserHandler) GetUserTransactionActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		id := c.Param("id")
		if id == "" {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}

		userID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		numSuccessfulOrders, numFailedOrders, totalOrders, err := h.service.GetUserTransactionActivity(userID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan data aktivitas customer: "+err.Error())
		}

		numSuccessfulChallenge, numFailedChallenge, totalChallenge, err := h.service.GetUserChallengeActivity(userID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan data aktivitas customer: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan data aktivitas customer",
			dto.FormatUserActivityResponse(
				numSuccessfulOrders, numFailedOrders, totalOrders,
				numSuccessfulChallenge, numFailedChallenge, totalChallenge,
			))
	}
}

func (h *UserHandler) GetUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		user, err := h.service.GetUsersById(currentUser.ID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail pengguna: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapat detail pengguna", dto.FormatUserProfileResponse(user))
	}
}
