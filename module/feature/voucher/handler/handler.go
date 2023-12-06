package handler

import (
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/dto"
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
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		voucherRequest := new(dto.CreateVoucherRequest)
		if err := c.Bind(voucherRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(voucherRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		newVoucher := &entities.VoucherModels{
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
			return response.SendStatusInternalServerResponse(c, "Gagal menambahkan kupon: "+err.Error())
		}

		return response.SendStatusCreatedResponse(c, "Kupon berhasil ditambahkan", dto.FormatVoucher(result))

	}
}

func (h *VoucherHandler) GetAllVouchers() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		status := c.QueryParam("status")
		category := c.QueryParam("category")

		var vouchers []*entities.VoucherModels
		var totalItems int64
		var err error

		if status != "" && category != "" {
			vouchers, totalItems, err = h.service.GetVoucherByStatusCategory(pageConv, perPage, status, category)
		} else if status != "" {
			vouchers, totalItems, err = h.service.GetVoucherByStatus(pageConv, perPage, status)
		} else if category != "" {
			vouchers, totalItems, err = h.service.GetVoucherByCategory(pageConv, perPage, category)
		} else {
			vouchers, totalItems, err = h.service.GetAllVoucher(pageConv, perPage)
		}

		if err != nil {
			c.Logger().Error("handler: failed to fetch all vouchers:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar kupon: "+err.Error())
		}
		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterVoucher(vouchers), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar kupon")
	}
}

func (h *VoucherHandler) UpdateVouchers() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.UpdateVoucherRequest)
		voucherID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")

		}
		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		updatedVoucher := &entities.VoucherModels{
			ID:          voucherID,
			Name:        req.Name,
			Code:        req.Code,
			Category:    req.Category,
			Description: req.Description,
			Discount:    req.Discount,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
			MinPurchase: req.MinPurchase,
			Stock:       req.Stock,
			Status:      req.Status,
		}
		err = h.service.UpdateVoucher(voucherID, updatedVoucher)
		if err != nil {
			c.Logger().Error("handler: failed create voucher:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui kupon: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil memperbarui kupon")
	}
}

func (h *VoucherHandler) DeleteVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		voucherID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}
		err = h.service.DeleteVoucher(voucherID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal menghapus kupon: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil menghapus kupon")
	}
}

func (h *VoucherHandler) GetVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		voucherID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}
		vouchers, err := h.service.GetVoucherById(voucherID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan voucher: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan voucher", dto.FormatVoucher(vouchers))
	}
}

func (h *VoucherHandler) ClaimVoucher() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.ClaimsVoucherRequest)
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		newClaims := &entities.VoucherClaimModels{
			UserID:    currentUser.ID,
			VoucherID: req.VoucherID,
		}
		if err := h.service.ClaimVoucher(newClaims); err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal klaim kupon: "+err.Error())
		}

		return response.SendStatusOkResponse(c, "Berhasil klaim kupon")

	}

}

func (h *VoucherHandler) GetVoucherUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		result, err := h.service.GetUserVouchers(currentUser.ID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan kupon user: "+err.Error())
		}

		formattedResponse, err := dto.GetVoucherUserFormatter(result)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memformat respons: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan kupon user", formattedResponse)
	}
}

func (h *VoucherHandler) GetAllVouchersToClaims() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		limit := 7
		result, err := h.service.GetAllVoucherToClaims(limit, currentUser.ID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar kupon: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan daftar kupon", dto.FormatterVoucherToClaims(result))
	}
}
