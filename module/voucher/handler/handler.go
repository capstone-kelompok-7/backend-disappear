package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
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
		var input = domain.VoucherModels{}
		c.Bind(&input)

		if input.Name == "" || input.Category == "" || input.Code == "" || input.Description == "" || input.Discouunt < 0 || input.EndDate == "" || input.StartDate == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid input",
			})
		}

		result, err := h.service.CreateVoucher(input)
		if err != nil {
			c.Logger().Error("handler: failed create voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error, %s", err.Error()))
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Voucher berhasil ditambahkan.",
			"data":    domain.VoucherResponseFormatter(*result),
		})
	}
}

func (h *VoucherHandler) GetAllVouchers() echo.HandlerFunc {
	return func(c echo.Context) error {
		page := c.QueryParam("page")
		search := c.QueryParam("search")
		pagee, err := strconv.Atoi(page)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Invalid page number")
		}
		limit := c.QueryParam("limit")
		limitt, err := strconv.Atoi(limit)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Invalid limit number")
		}

		if page == "" {
			page = "1"
		}

		if limit == "" {
			limit = "5"
		}

		prevPage := pagee - 1
		nextPage := pagee + 1

		if prevPage < 1 {
			prevPage = 1
		}

		result, err := h.service.GetAllVouchers(pagee, limitt, search)
		if err != nil {
			c.Logger().Error("handler: failed create voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error, %s", err.Error()))
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Berhasil.",
			"data":    domain.VoucherModelsFormatterAll(result),
			"pagination": map[string]interface{}{
				"current_page":  pagee,
				"toal_page":     len(result),
				"previous_page": prevPage,
				"next_page":     nextPage,
			},
		})
	}
}

func (h *VoucherHandler) EditVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = domain.VoucherModels{}
		var voucherid = c.Param("voucher_id")
		voucheridfix, _ := strconv.Atoi(voucherid)

		c.Bind(&input)

		if input.Name == "" || input.Category == "" || input.Code == "" || input.Description == "" || input.Discouunt < 0 || input.EndDate == "" || input.StartDate == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid input",
			})
		}

		input.ID = uint64(voucheridfix)
		result, err := h.service.EditVoucherById(input)
		if err != nil {
			c.Logger().Error("handler: failed edit voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error, %s", err.Error()))
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Voucher berhasil diperbarui.",
			"data":    domain.VoucherResponseFormatter(*result),
		})

	}
}
func (h *VoucherHandler) DeleteVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {

		var voucherid = c.Param("voucher_id")
		voucheridfix, _ := strconv.Atoi(voucherid)

		result := h.service.DeleteVoucherById(voucheridfix)
		if result != nil {
			c.Logger().Error("handler: failed delete voucher:", result.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error, %s", result.Error()))
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Voucher berhasil dihapus.",
		})

	}
}

func (h *VoucherHandler) GetVoucherById() echo.HandlerFunc {
	return func(c echo.Context) error {
		voucherid := c.Param("voucher_id")
		voucheridfix, _ := strconv.Atoi(voucherid)

		result, err := h.service.GetVoucherById(voucheridfix)
		if err != nil {
			c.Logger().Error("handler: failed get voucher:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error, %s", err.Error()))
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Berhasil.",
			"data":    domain.VoucherResponseFormatter(*result),
		})

	}
}
