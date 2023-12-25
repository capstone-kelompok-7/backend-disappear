package order

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order/dto"
	"github.com/labstack/echo/v4"
	"time"
)

type RepositoryOrderInterface interface {
	FindByName(page, perPage int, name string) ([]*entities.OrderModels, error)
	GetTotalCustomerCountByName(name string) (int64, error)
	FindAll(page, perPage int) ([]*entities.OrderModels, error)
	GetTotalOrderCount() (int64, error)
	GetOrderById(orderID string) (*entities.OrderModels, error)
	CreateOrder(newOrder *entities.OrderModels) (*entities.OrderModels, error)
	ConfirmPayment(orderID string, orderStatus, paymentStatus string) error
	ProcessGatewayPayment(totalAmountPaid uint64, orderID string, paymentMethod, name, email string) (interface{}, error)
	CheckTransaction(orderID string) (dto.Status, error)
	UpdateOrderStatus(req *dto.UpdateOrderStatus) error
	GetAllOrdersByUserID(userID uint64) ([]*entities.OrderModels, error)
	GetAllOrdersWithFilter(userID uint64, orderStatus string) ([]*entities.OrderModels, error)
	AcceptOrder(orderID, orderStatus string) error
	Tracking(courier, awb string) (map[string]interface{}, error)
	GetOrderByDateRange(startDate, endDate time.Time, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountByDateRange(startDate, endDate time.Time) (int64, error)
	GetOrderByOrderStatus(orderStatus string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountByByOrderStatus(orderStatus string) (int64, error)
	GetOrderByDateRangeAndStatus(startDate, endDate time.Time, status string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountByDateRangeAndStatus(startDate, endDate time.Time, status string) (int64, error)
	GetOrderByDateRangeAndStatusAndSearch(startDate, endDate time.Time, status, search string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountByDateRangeAndStatusAndSearch(startDate, endDate time.Time, status, search string) (int64, error)
	GetOrdersBySearchAndDateRange(startDate, endDate time.Time, search string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountBySearchAndDateRange(startDate, endDate time.Time, search string) (int64, error)
	GetOrdersBySearchAndStatus(status, search string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrdersCountBySearchAndStatus(status, search string) (int64, error)
	GetOrderBySearchAndPaymentStatus(status, search string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountBySearchAndPaymentStatus(status, search string) (int64, error)
	GetOrderByDateRangeAndPaymentStatusAndSearch(startDate, endDate time.Time, status, search string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountByDateRangeAndPaymentStatusAndSearch(startDate, endDate time.Time, status, search string) (int64, error)
	GetOrderByPaymentStatus(orderStatus string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountByByPaymentStatus(orderStatus string) (int64, error)
	GetOrderByDateRangeAndPaymentStatus(startDate, endDate time.Time, status string, offset, limit int) ([]*entities.OrderModels, error)
	GetOrderCountByDateRangeAndPaymentStatus(startDate, endDate time.Time, status string) (int64, error)
}

type ServiceOrderInterface interface {
	GetAll(page, perPage int) ([]*entities.OrderModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetOrdersByName(page, perPage int, name string) ([]*entities.OrderModels, int64, error)
	GetOrderById(orderID string) (*entities.OrderModels, error)
	CreateOrder(userID uint64, request *dto.CreateOrderRequest) (interface{}, error)
	ConfirmPayment(orderID string) error
	CreateOrderFromCart(userID uint64, request *dto.CreateOrderCartRequest) (interface{}, error)
	CancelPayment(orderID string) error
	CallBack(notifPayload map[string]any) error
	UpdateOrderStatus(req *dto.UpdateOrderStatus) error
	GetAllOrdersByUserID(userID uint64) ([]*entities.OrderModels, error)
	GetAllOrdersWithFilter(userID uint64, orderStatus string) ([]*entities.OrderModels, error)
	AcceptOrder(orderID string) error
	Tracking(courier, awb string) (map[string]interface{}, error)
	GetOrderByDateRange(filterType string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrderByOrderStatus(orderStatus string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrderByDateRangeAndStatus(filterType, status string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrderByDateRangeAndStatusAndSearch(filterType, status, search string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrderBySearchAndDateRange(filterType, search string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrdersBySearchAndStatus(status, search string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrdersBySearchAndPaymentStatus(status, search string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrderByDateRangeAndPaymentStatusAndSearch(filterType, status, search string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrderByDateRangeAndPaymentStatus(filterType, status string, page, perPage int) ([]*entities.OrderModels, int64, error)
	GetOrderByPaymentStatus(orderStatus string, page, perPage int) ([]*entities.OrderModels, int64, error)
	ProcessManualPayment(orderID string) (*entities.OrderModels, error)
	ProcessGatewayPayment(totalAmountPaid uint64, orderID string, paymentMethod, name, email string) (interface{}, error)
	SendNotificationOrder(request dto.SendNotificationOrderRequest) (string, error)
	SendNotificationPayment(request dto.SendNotificationPaymentRequest) (string, error)
}

type HandlerOrderInterface interface {
	GetOrderById() echo.HandlerFunc
	CreateOrder() echo.HandlerFunc
	GetAllOrders() echo.HandlerFunc
	ConfirmPayment() echo.HandlerFunc
	CreateOrderFromCart() echo.HandlerFunc
	CancelPayment() echo.HandlerFunc
	Callback() echo.HandlerFunc
	UpdateOrderStatus() echo.HandlerFunc
	GetAllOrderByUserID() echo.HandlerFunc
	AcceptOrder() echo.HandlerFunc
	Tracking() echo.HandlerFunc
	GetAllPayment() echo.HandlerFunc
}
