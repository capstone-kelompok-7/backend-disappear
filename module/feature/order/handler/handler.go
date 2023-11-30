package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
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
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
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
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar pesanan: ")
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterOrder(orders), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Daftar pesanan")
	}
}

func (h *OrderHandler) GetOrderById() echo.HandlerFunc {
	return func(c echo.Context) error {
		orderID := c.Param("id")
		if orderID == "" {
			return response.SendBadRequestResponse(c, "Format ID yang Anda masukkan tidak sesuai.")
		}
		result, err := h.service.GetOrderById(orderID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan detail pesanan: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan detail pesanan", dto.FormatOrderDetail(result))
	}
}

func (h *OrderHandler) CreateOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.CreateOrderRequest)
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		result, err := h.service.CreateOrder(currentUser.ID, req)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal membuat pesanan: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil membuat pesanan, silahkan melakukan konfirmasi pembayaran", dto.CreateOrderFormatter(result))
	}
}

func (h *OrderHandler) ConfirmPayment() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		orderID := c.Param("id")
		if orderID == "" {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}
		err := h.service.ConfirmPayment(orderID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mmelakukan pesanan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Pembayaran berhasil dikonfimasi, silahkan memproses pesanan")
	}
}

func (h *OrderHandler) CreateOrderFromCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		req := new(dto.CreateOrderCartRequest)
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}
		result, err := h.service.CreateOrderFromCart(currentUser.ID, req)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal membuat pesanan dari keranjang: "+err.Error())
		}
		return response.SendStatusCreatedResponse(c, "Berhasil membuat pesanan dari keranjang, silahkan melakukan konfirmasi pembayaran", dto.CreateOrderFormatter(result))
	}
}

func (h *OrderHandler) CancelPayment() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		orderID := c.Param("id")
		if orderID == "" {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}
		err := h.service.CancelPayment(orderID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mmelakukan pesanan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Pembayaran berhasil dibatalkan")
	}
}
