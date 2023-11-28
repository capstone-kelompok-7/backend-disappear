package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service order.ServiceOrderInterface
}

func NewOrderHandler(service order.ServiceOrderInterface) order.HandlerOrderInterface {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetAllOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusForbidden, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageConv, _ := strconv.Atoi(strconv.Itoa(page))
		perPage := 8

		var orders []*entities.OrderModels
		var totalItems int64
		var err error
		search := c.QueryParam("search")
		if search != "" {
			orders, totalItems, err = h.service.GetOrdersByName(page, perPage, search)
		} else {
			orders, totalItems, err = h.service.GetAll(pageConv, perPage)
		}
		if err != nil {
			c.Logger().Error("handler: failed to fetch all orders:", err.Error())
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan pesanan: ")
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.Pagination(c, dto.FormatterOrder(orders), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Daftar pesanan")
	}
}

func (h *OrderHandler) GetOrderById() echo.HandlerFunc {
	return func(c echo.Context) error {
		orderID := c.Param("id")
		if orderID == "" {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		result, err := h.service.GetOrderById(orderID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mendapatkan detail pesanan: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan detail pesanan", dto.FormatOrderDetail(result))
	}
}

func (h *OrderHandler) CreateOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendErrorResponse(c, http.StatusForbidden, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.CreateOrderRequest)
		if err := c.Bind(req); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Validasi gagal: "+err.Error())
		}
		result, err := h.service.CreateOrder(currentUser.ID, req)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal membuat pesanan: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil membuat pesanan, silahkan melakukan konfirmasi pembayaran", dto.CreateOrderFormatter(result))
	}
}

func (h *OrderHandler) ConfirmPayment() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendErrorResponse(c, http.StatusForbidden, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		orderID := c.Param("id")
		if orderID == "" {
			return response.SendErrorResponse(c, http.StatusBadRequest, "Format input yang Anda masukkan tidak sesuai.")
		}
		err := h.service.ConfirmPayment(orderID)
		if err != nil {
			return response.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mmelakukan pesanan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Pembayaran berhasil dikonfimasi, silahkan memproses pesanan")
	}
}
