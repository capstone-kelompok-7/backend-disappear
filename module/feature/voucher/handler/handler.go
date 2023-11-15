package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type VoucherHandler struct {
	service voucher.ServiceVoucherInterface
}

func NewVoucherHandler(service voucher.ServiceVoucherInterface) voucher.HandlerVoucherInterface {
	return &VoucherHandler{
		service: service,
	}
}

func (h *VoucherHandler) CreateVoucher() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		voucherRequest := new(dto.CreateVoucherRequest)
		if err := c.Bind(&voucherRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(voucherRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		newVoucher := entities.VoucherModels{
			Name:        voucherRequest.Name,
			Code:        voucherRequest.Code,
			Category:    voucherRequest.Category,
			Description: voucherRequest.Description,
			Discount:    voucherRequest.Discount,
			StartDate:   voucherRequest.StartDate,
			EndDate:     voucherRequest.EndDate,
			MinPurchase: voucherRequest.MinPurchase,
			Stock:       voucherRequest.Stock,
		}

		result, err := h.service.CreateVoucher(newVoucher)
		if err != nil {
			c.Logger().Error("handler: failed create voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Kupon berhasil ditambahkan", dto.FormatVoucher(result))

	}
}

func (h *VoucherHandler) GetAllVouchers() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var vouchers []entities.VoucherModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			vouchers, totalItems, err = h.service.GetVouchersByName(page, perPage, search)
		} else {
			vouchers, totalItems, err = h.service.GetAllVoucher(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all vouchers:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
		}

		current_page, total_pages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(current_page, total_pages)
		prevPage := h.service.GetPrevPage(current_page)

		return response.Pagination(c, dto.FormatterVoucher(vouchers), current_page, total_pages, int(totalItems), nextPage, prevPage, "Daftar kupon")
	}
}

func (h *VoucherHandler) UpdateVouchers() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		var voucherRequest dto.UpdateVoucherRequest
		voucherID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		if err := c.Bind(&voucherRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(&voucherRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		_, err = h.service.UpdateVoucher(voucherID, voucherRequest)
		if err != nil {
			c.Logger().Error("handler: failed create voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Kesalahan Server Internal: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Kupon berhasil dirubah")
	}
}

func (h *VoucherHandler) DeleteVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusUnauthorized, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		voucherID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteVoucher(voucherID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus kupon: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil hapus kupon")
	}
}

func (h *VoucherHandler) GetVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		voucherID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		vouchers, err := h.service.GetVoucherById(voucherID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan voucher: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan voucher", dto.FormatVoucher(vouchers))
	}
}
