package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/capstone-kelompok-7/backend-disappear/utils/upload"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"net/http"
	"strconv"
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
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		email := c.QueryParam("email")
		if email == "" {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Parameter Email tidak ada")
		}
		user, err := h.service.GetUsersByEmail(email)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		}
		return response.SendSuccessResponse(c, "Berhasil mendapat data pengguna", user)
	}
}

func (h *UserHandler) ChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var updateRequest dto.UpdatePasswordRequest
		if err := c.Bind(&updateRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		if err := utils.ValidateStruct(updateRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		currentUser := c.Get("CurrentUser").(*entities.UserModels)

		err := h.service.ValidatePassword(currentUser.ID, updateRequest.OldPassword, updateRequest.NewPassword, updateRequest.ConfirmPassword)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui kata sandi: "+err.Error())
		}
		err = h.service.ChangePassword(currentUser.ID, updateRequest)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Kata sandi berhasil diperbarui")
	}
}

func (h *UserHandler) GetUsersById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		id := c.Param("id")
		if id == "" {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Parameter id tidak ada")
		}

		userID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		user, err := h.service.GetUsersById(userID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		}

		return response.SendSuccessResponse(c, "Berhasil mendapat data pengguna", user)
	}
}

func (h *UserHandler) GetAllUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var users []*entities.UserModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			users, totalItems, err = h.service.GetUsersByName(page, perPage, search)
			if err != nil {
				c.Logger().Error("handler: failed to fetch users by name:", err.Error())
				return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
			}
		} else {
			users, totalItems, err = h.service.GetAllUsers(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all users:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, dto.FormatterUsers(users), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar customer")
	}
}

func (h *UserHandler) EditProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		var editProfileRequest dto.EditProfileRequest
		if err := c.Bind(&editProfileRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		if err := utils.ValidateStruct(editProfileRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
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
			editProfileRequest.PhotoProfile = uploadedURL
		}
		_, err = h.service.EditProfile(currentUser.ID, editProfileRequest)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui profil: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Profil berhasil diperbarui")
	}
}

func (h *UserHandler) DeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteAccount(userID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus pengguna: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil hapus pengguna")
	}
}
