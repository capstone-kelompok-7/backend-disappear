package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"net/http"
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
		addressRequest := new(dto.CreateAddressRequest)
		if err := c.Bind(addressRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(addressRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		newAddress := &entities.AddressModels{
			UserID:       currentUser.ID,
			AcceptedName: addressRequest.AcceptedName,
			Street:       addressRequest.Street,
			SubDistrict:  addressRequest.SubDistrict,
			City:         addressRequest.City,
			Province:     addressRequest.Province,
			PostalCode:   addressRequest.PostalCode,
			Note:         addressRequest.Note,
		}
		createdAddress, err := h.service.CreateAddress(newAddress)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil menambahkan alamat", dto.FormatAddress(createdAddress))
	}
}

func (h *AddressHandler) GetAllAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8
		var addresses []*entities.AddressModels
		var totalItems int64
		var err error
		addresses, totalItems, err = h.service.GetAll(currentUser.ID, pageConv, perPage)
		if err != nil {
			c.Logger().Error("handler: failed to fetch all addresses:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.Pagination(c, dto.FormatterAddress(addresses), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Daftar alamat")
	}
}