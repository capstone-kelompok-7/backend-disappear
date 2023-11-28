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
}

type ServiceOrderInterface interface {
	GetAll(page, perPage int) ([]*entities.OrderModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetOrdersByName(page, perPage int, name string) ([]*entities.OrderModels, int64, error)
	GetOrderById(orderID string) (*entities.OrderModels, error)
	CreateOrder(userID uint64, request *dto.CreateOrderRequest) (*entities.OrderModels, error)
	ConfirmPayment(orderID string) error
}

type HandlerOrderInterface interface {
	GetOrderById() echo.HandlerFunc
	CreateOrder() echo.HandlerFunc
	GetAllOrders() echo.HandlerFunc
	ConfirmPayment() echo.HandlerFunc
}
