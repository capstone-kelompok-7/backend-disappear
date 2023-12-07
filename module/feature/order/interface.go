package order

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order/dto"
	"github.com/labstack/echo/v4"
)

type RepositoryOrderInterface interface {
	FindByName(page, perPage int, name string) ([]*entities.OrderModels, error)
	GetTotalCustomerCountByName(name string) (int64, error)
	FindAll(page, perPage int) ([]*entities.OrderModels, error)
	GetTotalOrderCount() (int64, error)
	GetOrderById(orderID string) (*entities.OrderModels, error)
	CreateOrder(newOrder *entities.OrderModels) (*entities.OrderModels, error)
	ConfirmPayment(orderID string, orderStatus, paymentStatus string) error
	ProcessGatewayPayment(totalAmountPaid uint64, orderID string, paymentMethod string) (interface{}, error)
	CheckTransaction(orderID string) (dto.Status, error)
	UpdateOrderStatus(req *dto.UpdateOrderStatus) error
	GetAllOrdersByUserID(userID uint64) ([]*entities.OrderModels, error)
	GetAllOrdersWithFilter(userID uint64, orderStatus string) ([]*entities.OrderModels, error)
	AcceptOrder(orderID, orderStatus string) error
	Tracking(courier, awb string) (map[string]interface{}, error)
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
}
