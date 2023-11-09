package handler

import (
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
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
		var voucherRequest = new(domain.VoucherRequestModel)

		if err := c.Bind(&voucherRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(voucherRequest); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}

		newVoucher := &domain.VoucherModels{
			ID:          voucherRequest.ID,
			Name:        voucherRequest.Name,
			Code:        voucherRequest.Code,
			Category:    voucherRequest.Category,
			Description: voucherRequest.Description,
			Discouunt:   voucherRequest.Discouunt,
			StartDate:   voucherRequest.StartDate,
			EndDate:     voucherRequest.EndDate,
			MinAmount:   voucherRequest.MinAmount,
		}

		result, err := h.service.CreateVoucher(*newVoucher)
		if err != nil {
			c.Logger().Error("handler: failed create voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error, "+err.Error())
		}

		return response.SendSuccessResponse(c, "Voucher berhasil ditambahkan", domain.VoucherResponseFormatterCreate(result))
	}
}

func (h *VoucherHandler) GetAllVouchers() echo.HandlerFunc {
	return func(c echo.Context) error {
		page := c.QueryParam("page")
		category := c.QueryParam("category")
		pageconv, _ := strconv.Atoi(page)

		prevPage := pageconv - 1
		nextPage := pageconv + 1

		allvoucher, _ := h.service.GetAllVouchersToCalculatePage()
		var calculatePage = len(allvoucher) / 8

		if prevPage < 1 {
			prevPage = 1
		}

		if pageconv == 0 {
			pageconv = 1
		}

		result, err := h.service.GetAllVouchers(pageconv, 8, category)
		if err != nil {
			c.Logger().Error("handler: failed create voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error, "+err.Error())
		}
		return response.Pagination(c, domain.VoucherModelsFormatterAll(result), pageconv, calculatePage, len(result), nextPage, prevPage, "Berhasil")
	}
}

func (h *VoucherHandler) EditVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = domain.VoucherModels{}
		var voucherid = c.Param("voucher_id")
		voucheridfix, _ := strconv.Atoi(voucherid)

		if err := c.Bind(&input); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if input.Name == "" || input.Code == "" || input.Category == "" || input.Description == "" || input.Discouunt < 0 || input.MinAmount < 0 || input.StartDate == "" || input.EndDate == "" {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input")
		}

		input.ID = uint64(voucheridfix)
		result, err := h.service.EditVoucherById(input)
		if err != nil {
			c.Logger().Error("handler: failed edit voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error, "+err.Error())
		}
		return response.SendSuccessResponse(c, "Voucher berhasil diperbarui", domain.VoucherResponseFormatterCreate(result))
	}
}
func (h *VoucherHandler) DeleteVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var voucherid = c.Param("voucher_id")
		voucheridfix, _ := strconv.Atoi(voucherid)

		result := h.service.DeleteVoucherById(voucheridfix)
		if result != nil {
			c.Logger().Error("handler: failed delete voucher:", result.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}
		return response.SendStatusOkResponse(c, "Voucher berhasil dihapus.")
	}
}

func (h *VoucherHandler) GetVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		voucherid := c.Param("voucher_id")
		voucherconv, _ := strconv.Atoi(voucherid)

		result, err := h.service.GetVoucherById(voucherconv)
		if err != nil {
			c.Logger().Error("handler: failed get voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error, "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil.", domain.VoucherResponseFormatter(result))
	}
}
