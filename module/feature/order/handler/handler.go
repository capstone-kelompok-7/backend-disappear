package handler

import (
	"encoding/json"
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
		dateFilter := c.QueryParam("date_filter")
		statusFilter := c.QueryParam("status_filter")

		if search != "" && dateFilter != "" && statusFilter != "" {
			orders, totalItems, err = h.service.GetOrderByDateRangeAndStatusAndSearch(dateFilter, statusFilter, search, page, perPage)
		} else if search != "" && statusFilter != "" {
			orders, totalItems, err = h.service.GetOrdersBySearchAndStatus(statusFilter, search, page, perPage)
		} else if search != "" && dateFilter != "" {
			orders, totalItems, err = h.service.GetOrderBySearchAndDateRange(dateFilter, search, page, perPage)
		} else if dateFilter != "" && statusFilter != "" {
			orders, totalItems, err = h.service.GetOrderByDateRangeAndStatus(dateFilter, statusFilter, page, perPage)
		} else if search != "" {
			orders, totalItems, err = h.service.GetOrdersByName(page, perPage, search)
		} else if dateFilter != "" {
			orders, totalItems, err = h.service.GetOrderByDateRange(dateFilter, page, perPage)
		} else if statusFilter != "" {
			orders, totalItems, err = h.service.GetOrderByOrderStatus(statusFilter, page, perPage)
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

		return response.SendPaginationResponse(c, dto.FormatterOrder(orders), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar pesanan")
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
		switch req.PaymentMethod {
		case "whatsapp", "telegram":
			orderResult := result.(*entities.OrderModels)
			return response.SendStatusCreatedResponse(c, "Berhasil membuat pesanan, silahkan melakukan konfirmasi pembayaran", dto.CreateOrderFormatter(orderResult))
		case "qris", "bank_transfer", "mandiri", "bca", "gopay", "shopepay":
			return response.SendStatusCreatedResponse(c, "Berhasil membuat pesanan dengan payment gateway", result)
		default:
			return response.SendBadRequestResponse(c, "Metode pembayaran tidak valid")
		}
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
		switch req.PaymentMethod {
		case "whatsapp", "telegram":
			orderResult := result.(*entities.OrderModels)
			return response.SendStatusCreatedResponse(c, "Berhasil membuat pesanan, silahkan melakukan konfirmasi pembayaran", dto.CreateOrderFormatter(orderResult))
		case "qris", "bank_transfer", "gopay":
			return response.SendStatusCreatedResponse(c, "Berhasil membuat pesanan dengan payment gateway", result)
		default:
			return response.SendBadRequestResponse(c, "Metode pembayaran tidak valid")
		}
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

func (h *OrderHandler) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		var notificationPayload map[string]any

		if err := json.NewDecoder(c.Request().Body).Decode(&notificationPayload); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		err := h.service.CallBack(notificationPayload)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal callback pesanan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil callback")
	}
}

func (h *OrderHandler) UpdateOrderStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "admin" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		req := new(dto.UpdateOrderStatus)
		if err := c.Bind(req); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}

		if err := utils.ValidateStruct(req); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		if err := h.service.UpdateOrderStatus(req); err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal memperbarui status pesanan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil memperbarui status pesanan")
	}
}

func (h *OrderHandler) GetAllOrderByUserID() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		orderStatus := c.QueryParam("order_status")

		var result []*entities.OrderModels
		var err error

		if orderStatus == "" {
			result, err = h.service.GetAllOrdersByUserID(currentUser.ID)
		} else {
			result, err = h.service.GetAllOrdersWithFilter(currentUser.ID, orderStatus)
		}

		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan data order customer: "+err.Error())
		}
		return response.SendSuccessResponse(c, "Berhasil mendapatkan data order customer", dto.FormatterGetAllOrderUser(result))
	}
}

func (h *OrderHandler) AcceptOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}
		orderID := c.Param("id")
		if orderID == "" {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai")
		}
		err := h.service.AcceptOrder(orderID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mengkonfirmasi pesanan: "+err.Error())
		}
		return response.SendStatusOkResponse(c, "Berhasil mengkonfirmasi pesanan")
	}
}

func (h *OrderHandler) Tracking() echo.HandlerFunc {
	return func(c echo.Context) error {
		courier := c.QueryParam("courier")
		awb := c.QueryParam("awb")
		result, err := h.service.Tracking(courier, awb)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan resi: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan resi", result)
	}
}

func (h *OrderHandler) GetAllPayment() echo.HandlerFunc {
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
		dateFilter := c.QueryParam("date_filter")
		statusFilter := c.QueryParam("status_filter")

		if search != "" && dateFilter != "" && statusFilter != "" {
			orders, totalItems, err = h.service.GetOrderByDateRangeAndPaymentStatusAndSearch(dateFilter, statusFilter, search, page, perPage)
		} else if search != "" && statusFilter != "" {
			orders, totalItems, err = h.service.GetOrdersBySearchAndPaymentStatus(statusFilter, search, page, perPage)
		} else if search != "" && dateFilter != "" {
			orders, totalItems, err = h.service.GetOrderBySearchAndDateRange(dateFilter, search, page, perPage)
		} else if dateFilter != "" && statusFilter != "" {
			orders, totalItems, err = h.service.GetOrderByDateRangeAndPaymentStatus(dateFilter, statusFilter, page, perPage)
		} else if search != "" {
			orders, totalItems, err = h.service.GetOrdersByName(page, perPage, search)
		} else if dateFilter != "" {
			orders, totalItems, err = h.service.GetOrderByDateRange(dateFilter, page, perPage)
		} else if statusFilter != "" {
			orders, totalItems, err = h.service.GetOrderByPaymentStatus(statusFilter, page, perPage)
		} else {
			orders, totalItems, err = h.service.GetAll(pageConv, perPage)
		}

		if err != nil {
			c.Logger().Error("handler: failed to fetch all orders:", err.Error())
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar pembayaran: ")
		}

		currentPage, totalPages := h.service.CalculatePaginationValues(pageConv, int(totalItems), perPage)
		nextPage := h.service.GetNextPage(currentPage, totalPages)
		prevPage := h.service.GetPrevPage(currentPage)

		return response.SendPaginationResponse(c, dto.FormatterOrderPayment(orders), currentPage, totalPages, int(totalItems), nextPage, prevPage, "Berhasil mendapatkan daftar pembayaran")
	}
}
