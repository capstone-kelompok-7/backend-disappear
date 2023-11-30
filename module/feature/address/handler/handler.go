package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"strconv"
)

type AddressHandler struct {
	service address.ServiceAddressInterface
}

func NewAddressHandler(service address.ServiceAddressInterface) address.HandlerAddressInterface {
	return &AddressHandler{
		service: service,
	}
}

func (h *AddressHandler) CreateAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		addressRequest := new(dto.CreateAddressRequest)
		if err := c.Bind(addressRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(addressRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		newAddress := &entities.AddressModels{
			UserID:       currentUser.ID,
			AcceptedName: addressRequest.AcceptedName,
			Phone:        addressRequest.Phone,
			Address:      addressRequest.Address,
			IsPrimary:    addressRequest.IsPrimary,
		}
		createdAddress, err := h.service.CreateAddress(newAddress)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan alamat: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil menambahkan alamat", dto.FormatAddress(createdAddress))
	}
}

func (h *AddressHandler) GetAllAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8
		var addresses []*entities.AddressModels
		var totalItems int64
		var err error
		addresses, totalItems, err = h.service.GetAll(currentUser.ID, pageConv, perPage)
		if err != nil {
			c.Logger().Error("handler: failed to fetch all addresses:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan alamat: "+err.Error())
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterAddress(addresses), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar alamat")
	}
}

func (h *AddressHandler) UpdateAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		addressID := c.Param("id")
		id, err := strconv.ParseUint(addressID, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}
		var updatedAddress *dto.UpdateAddressRequest
		if err := c.Bind(&updatedAddress); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}
		if err := utils.ValidateStruct(updatedAddress); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		err = h.service.UpdateAddress(currentUser.ID, id, updatedAddress)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui alamat: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil memperbarui alamat")
	}
}

func (h *AddressHandler) DeleteAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		id := c.Param("id")
		addressID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai")
		}

		err = h.service.DeleteAddress(addressID, currentUser.ID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus alamat: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil menghapus alamat")
	}
}
