package service

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	addressMock "github.com/capstone-kelompok-7/backend-disappear/module/feature/address/mocks"
	address "github.com/capstone-kelompok-7/backend-disappear/module/feature/address/service"
	assistantMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/mocks"
	assistants "github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/service"
	cartMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/cart/mocks"
	cart "github.com/capstone-kelompok-7/backend-disappear/module/feature/cart/service"
	fcmMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm/mocks"
	fcm "github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm/service"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order/dto"
	orders "github.com/capstone-kelompok-7/backend-disappear/module/feature/order/mocks"
	productsMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/product/mocks"
	products "github.com/capstone-kelompok-7/backend-disappear/module/feature/product/service"
	userMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/mocks"
	user "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/service"
	voucherMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/mocks"
	vouchers "github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/service"
	utils "github.com/capstone-kelompok-7/backend-disappear/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupOrderService(t *testing.T) (
	*OrderService,
	*orders.RepositoryOrderInterface,
	*userMocks.RepositoryUserInterface,
	*productsMocks.RepositoryProductInterface,
	*assistantMocks.RepositoryAssistantInterface,
	*voucherMocks.RepositoryVoucherInterface,
	*addressMock.RepositoryAddressInterface,
	*cartMocks.RepositoryCartInterface,
	*fcmMocks.RepositoryFcmInterface,
	*utils.GeneratorInterface,
) {
	orderRepo := orders.NewRepositoryOrderInterface(t)
	generatorRepo := utils.NewGeneratorInterface(t)
	hashRepo := utils.NewHashInterface(t)
	productRepo := productsMocks.NewRepositoryProductInterface(t)
	assistantRepo := assistantMocks.NewRepositoryAssistantInterface(t)
	userRepo := userMocks.NewRepositoryUserInterface(t)
	voucherRepo := voucherMocks.NewRepositoryVoucherInterface(t)
	addressRepo := addressMock.NewRepositoryAddressInterface(t)
	cartRepo := cartMocks.NewRepositoryCartInterface(t)
	fcmRepo := fcmMocks.NewRepositoryFcmInterface(t)

	assistantService := assistants.NewAssistantService(assistantRepo, nil, config.Config{})
	productService := products.NewProductService(productRepo, assistantService)
	userService := user.NewUserService(userRepo, hashRepo)
	voucherService := vouchers.NewVoucherService(voucherRepo, userService)
	addressService := address.NewAddressService(addressRepo)
	cartService := cart.NewCartService(cartRepo, productService)
	fcmService := fcm.NewFcmService(fcmRepo)
	orderService := NewOrderService(orderRepo, generatorRepo, productService, voucherService, addressService, userService, cartService, fcmService)

	return orderService.(*OrderService), orderRepo, userRepo, productRepo, assistantRepo, voucherRepo, addressRepo, cartRepo, fcmRepo, generatorRepo
}

func TestGetFilterDateRange(t *testing.T) {
	service := &OrderService{}

	tests := []struct {
		filterType    string
		expectedStart time.Time
		expectedEnd   time.Time
		expectedError error
	}{
		{
			filterType:    "minggu ini",
			expectedStart: time.Now().In(time.UTC).AddDate(0, 0, -int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour),
			expectedEnd:   time.Now().In(time.UTC).AddDate(0, 0, 7-int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour),
			expectedError: nil,
		},
		{
			filterType:    "bulan ini",
			expectedStart: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
			expectedError: nil,
		},
		{
			filterType:    "tahun ini",
			expectedStart: time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(time.Now().Year()+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
			expectedError: nil,
		},
		{
			filterType:    "invalid",
			expectedError: errors.New("tipe filter tidak valid"),
		},
	}

	for _, test := range tests {
		startDate, endDate, err := service.GetFilterDateRange(test.filterType)

		if !startDate.Equal(test.expectedStart) || !endDate.Equal(test.expectedEnd) || !reflect.DeepEqual(err, test.expectedError) {
			t.Errorf("For filter type %s, expected (%v, %v, %v), but got (%v, %v, %v)", test.filterType, test.expectedStart, test.expectedEnd, test.expectedError, startDate, endDate, err)
		}
	}
}

func TestOrderService_PaginationFunctions(t *testing.T) {
	orderService := &OrderService{}

	t.Run("CalculatePaginationValues", func(t *testing.T) {
		// Test case 1
		pageInt, totalPages := orderService.CalculatePaginationValues(1, 20, 5)
		assert.Equal(t, 1, pageInt)
		assert.Equal(t, 4, totalPages)

		// Test case 2
		pageInt, totalPages = orderService.CalculatePaginationValues(-1, 15, 5)
		assert.Equal(t, 1, pageInt)
		assert.Equal(t, 3, totalPages)

		// Test case 3
		pageInt, totalPages = orderService.CalculatePaginationValues(7, 50, 10)
		assert.Equal(t, 5, pageInt)
		assert.Equal(t, 5, totalPages)
	})

	t.Run("GetNextPage", func(t *testing.T) {
		// Test case 1
		nextPage := orderService.GetNextPage(3, 5)
		assert.Equal(t, 4, nextPage)

		// Test case 2
		nextPage = orderService.GetNextPage(5, 5)
		assert.Equal(t, 5, nextPage)

		// Test case 3
		nextPage = orderService.GetNextPage(8, 10)
		assert.Equal(t, 9, nextPage)
	})

	t.Run("GetPrevPage", func(t *testing.T) {
		// Test case 1
		prevPage := orderService.GetPrevPage(3)
		assert.Equal(t, 2, prevPage)

		// Test case 2
		prevPage = orderService.GetPrevPage(1)
		assert.Equal(t, 1, prevPage)

		// Test case 3
		prevPage = orderService.GetPrevPage(7)
		assert.Equal(t, 6, prevPage)
	})
}

func TestOrderService_GetAll(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)

	order := []*entities.OrderModels{
		{ID: "order1100sada", IdOrder: "13123sasdasd", AddressID: 1, UserID: 1},
		{ID: "123123order", IdOrder: "123888jjjhhss", AddressID: 1, UserID: 1},
	}

	t.Run("Success Case - Order Found", func(t *testing.T) {
		expectedTotalItems := int64(10)
		orderRepo.On("FindAll", 1, 8).Return(order, nil).Once()
		orderRepo.On("GetTotalOrderCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetAll(1, 8)

		assert.NoError(t, err)
		assert.Equal(t, len(order), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetTotalOrderCount Error", func(t *testing.T) {
		expectedErr := errors.New(" GetTotalOrderCount Error")
		orderRepo.On("FindAll", 1, 8).Return(order, nil).Once()
		orderRepo.On("GetTotalOrderCount").Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetAll(1, 8)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetAll Error", func(t *testing.T) {
		expectedErr := errors.New("FindAll Error")
		orderRepo.On("FindAll", 1, 8).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetAll(1, 8)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		orderRepo.AssertExpectations(t)
	})

}

func TestOrderService_GetOrdersByName(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)

	order := []*entities.OrderModels{
		{ID: "order1100sada", IdOrder: "13123sasdasd", AddressID: 1, UserID: 1},
		{ID: "123123order", IdOrder: "123888jjjhhss", AddressID: 1, UserID: 1},
	}

	name := "Test"

	t.Run("Success Case - Orders Found by Name", func(t *testing.T) {
		expectedTotalItems := int64(10)
		orderRepo.On("FindByName", 1, 8, name).Return(order, nil).Once()
		orderRepo.On("GetTotalCustomerCountByName", name).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetOrdersByName(1, 8, name)

		assert.NoError(t, err)
		assert.Equal(t, len(order), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Orders by Name", func(t *testing.T) {
		expectedErr := errors.New("failed to find orders by name")
		orderRepo.On("FindByName", 1, 8, name).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrdersByName(1, 8, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Order Count by Name", func(t *testing.T) {
		expectedErr := errors.New("failed to get total order count by name")
		orderRepo.On("FindByName", 1, 8, name).Return(order, nil).Once()
		orderRepo.On("GetTotalCustomerCountByName", name).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrdersByName(1, 8, name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrderById(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)

	order := &entities.OrderModels{
		ID:        "order1100sada",
		IdOrder:   "13123sasdasd",
		AddressID: 1,
		UserID:    1,
	}

	expectedOrder := &entities.OrderModels{
		ID:        order.ID,
		IdOrder:   order.IdOrder,
		AddressID: order.AddressID,
		UserID:    order.UserID,
	}

	t.Run("Success Case - Order Found", func(t *testing.T) {
		orderID := "order1100sada"
		orderRepo.On("GetOrderById", orderID).Return(expectedOrder, nil).Once()

		result, err := orderService.GetOrderById(orderID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedOrder.ID, result.ID)
		assert.Equal(t, expectedOrder.IdOrder, result.IdOrder)
		assert.Equal(t, expectedOrder.AddressID, result.AddressID)
		assert.Equal(t, expectedOrder.UserID, result.UserID)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {
		orderID := "unknown_order_id"

		expectedErr := errors.New("gagal mendapatkan pesanan")
		orderRepo.On("GetOrderById", orderID).Return(nil, expectedErr).Once()

		result, err := orderService.GetOrderById(orderID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_ProcessManualPayment(t *testing.T) {

	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	orderID := "order123"
	expectedOrder := &entities.OrderModels{
		ID:        orderID,
		IdOrder:   "123order",
		AddressID: 1,
		UserID:    1,
	}

	t.Run("Success Case - Process Manual Payment", func(t *testing.T) {

		orderRepo.On("GetOrderById", orderID).Return(expectedOrder, nil).Once()

		result, err := orderService.ProcessManualPayment(orderID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedOrder.ID, result.ID)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {
		expectedErr := errors.New("order not found")
		orderRepo.On("GetOrderById", orderID).Return(nil, expectedErr).Once()

		result, err := orderService.ProcessManualPayment(orderID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		orderRepo.AssertExpectations(t)
	})

}

func TestOrderService_ProcessGatewayPayment(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	orderID := "order123"
	totalAmountPaid := uint64(50000)
	paymentMethod := "credit_card"

	t.Run("Success Case - Process Gateway Payment", func(t *testing.T) {
		expectedResult := map[string]interface{}{
			"payment_status": "success",
		}

		orderRepo.On("ProcessGatewayPayment", totalAmountPaid, orderID, paymentMethod).Return(expectedResult, nil).Once()

		result, err := orderService.ProcessGatewayPayment(totalAmountPaid, orderID, paymentMethod)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedResult, result)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Payment Failure", func(t *testing.T) {
		expectedErr := errors.New("payment failed")
		orderRepo.On("ProcessGatewayPayment", totalAmountPaid, orderID, paymentMethod).Return(nil, expectedErr).Once()

		result, err := orderService.ProcessGatewayPayment(totalAmountPaid, orderID, paymentMethod)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		orderRepo.AssertExpectations(t)
	})

}

func TestOrderService_GetAllOrdersByUserID(t *testing.T) {
	orderService, orderRepo, userRepo, _, _, _, _, _, _, _ := setupOrderService(t)

	userID := uint64(123)
	expectedUser := &entities.UserModels{ID: userID}
	expectedOrders := []*entities.OrderModels{
		{ID: "order123", UserID: userID},
		{ID: "order321", UserID: 221}}

	t.Run("Success Case - Get Orders by User ID", func(t *testing.T) {

		userRepo.On("GetUsersById", userID).Return(expectedUser, nil).Once()
		orderRepo.On("GetAllOrdersByUserID", userID).Return(expectedOrders, nil).Once()

		result, err := orderService.GetAllOrdersByUserID(userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedUser.ID, result[0].UserID)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		userRepo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		result, err := orderService.GetAllOrdersByUserID(userID)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error GetAllOrdersByUserID", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		userRepo.On("GetUsersById", userID).Return(expectedUser, nil).Once()
		orderRepo.On("GetAllOrdersByUserID", userID).Return(nil, expectedErr).Once()

		result, err := orderService.GetAllOrdersByUserID(userID)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		userRepo.AssertExpectations(t)
	})

}

func TestOrderService_GetAllOrdersWithFilter(t *testing.T) {
	orderService, orderRepo, userRepo, _, _, _, _, _, _, _ := setupOrderService(t)

	userID := uint64(123)
	orderStatus := "completed"
	expectedOrders := []*entities.OrderModels{
		{ID: "order123", UserID: userID, OrderStatus: orderStatus},
	}

	t.Run("Success Case - Get Orders with Filter", func(t *testing.T) {
		expectedUser := &entities.UserModels{ID: userID}
		userRepo.On("GetUsersById", userID).Return(expectedUser, nil).Once()

		orderRepo.On("GetAllOrdersWithFilter", userID, orderStatus).Return(expectedOrders, nil).Once()

		result, err := orderService.GetAllOrdersWithFilter(userID, orderStatus)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedUser.ID, result[0].UserID)
		assert.Equal(t, orderStatus, result[0].OrderStatus)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		userRepo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		result, err := orderService.GetAllOrdersWithFilter(userID, orderStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Orders", func(t *testing.T) {
		expectedUser := &entities.UserModels{ID: userID}
		userRepo.On("GetUsersById", userID).Return(expectedUser, nil).Once()

		expectedErr := errors.New("failed to fetch orders")
		orderRepo.On("GetAllOrdersWithFilter", userID, orderStatus).Return(nil, expectedErr).Once()

		result, err := orderService.GetAllOrdersWithFilter(userID, orderStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_Tracking(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)

	courier := "JNE"
	awb := "123456789"

	t.Run("Success Case - Tracking Info Found", func(t *testing.T) {
		expectedTrackingInfo := map[string]interface{}{
			"status":      "In Transit",
			"destination": "Your Location",
		}
		orderRepo.On("Tracking", courier, awb).Return(expectedTrackingInfo, nil).Once()

		result, err := orderService.Tracking(courier, awb)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTrackingInfo, result)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Tracking Info Not Found", func(t *testing.T) {
		expectedErr := errors.New("tracking information not found")
		orderRepo.On("Tracking", courier, awb).Return(nil, expectedErr).Once()

		result, err := orderService.Tracking(courier, awb)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrderByDateRange(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	userID := uint64(123)
	expectedOrders := []*entities.OrderModels{
		{ID: "order123", UserID: userID},
	}

	t.Run("Success Case - Filter by Week", func(t *testing.T) {
		filterType := "Minggu Ini"
		page := 1
		perPage := 8

		startOfWeek := time.Now().In(time.UTC).AddDate(0, 0, -int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)
		endOfWeek := time.Now().In(time.UTC).AddDate(0, 0, 7-int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)

		expectedTotalItems := int64(len(expectedOrders))
		orderRepo.On("GetOrderByDateRange", startOfWeek, endOfWeek, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRange", startOfWeek, endOfWeek).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetOrderByDateRange(filterType, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Success Case - Filter by Month", func(t *testing.T) {
		filterType := "Bulan Ini"
		page := 1
		perPage := 8

		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		nextMonth := startOfMonth.AddDate(0, 1, 0)
		endOfMonth := nextMonth.Add(-time.Second)

		expectedTotalItems := int64(len(expectedOrders))
		orderRepo.On("GetOrderByDateRange", startOfMonth, endOfMonth, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRange", startOfMonth, endOfMonth).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetOrderByDateRange(filterType, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Success Case - Filter by Year", func(t *testing.T) {
		filterType := "Tahun Ini"
		page := 1
		perPage := 8

		now := time.Now()
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		nextYear := startOfYear.AddDate(1, 0, 0)
		endOfYear := nextYear.Add(-time.Second)

		expectedTotalItems := int64(len(expectedOrders))
		orderRepo.On("GetOrderByDateRange", startOfYear, endOfYear, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRange", startOfYear, endOfYear).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetOrderByDateRange(filterType, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Invalid Filter Type", func(t *testing.T) {
		filterType := "Invalid Type"
		page := 1
		perPage := 8

		expectedErr := errors.New("tipe filter tidak valid")

		result, totalItems, err := orderService.GetOrderByDateRange(filterType, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr.Error(), err.Error())

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Fetching Orders", func(t *testing.T) {
		filterType := "Minggu Ini"
		page := 1
		perPage := 8

		expectedErr := errors.New("pesanan tidak ditemukan")

		orderRepo.On("GetOrderByDateRange", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRange(filterType, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Fetching Order Count", func(t *testing.T) {
		filterType := "Minggu Ini"
		page := 1
		perPage := 8

		expectedErr := errors.New("gagal mendapatkan total pesanan")

		orderRepo.On("GetOrderByDateRange", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRange", mock.Anything, mock.Anything).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRange(filterType, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

}

func TestOrderService_GetOrderByOrderStatus(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	expectedOrders := []*entities.OrderModels{
		{ID: "order123", OrderStatus: "Pending"},
		{ID: "order456", OrderStatus: "Shipped"},
	}

	t.Run("Success Case - Get Orders by Order Status", func(t *testing.T) {
		orderStatus := "Pending"
		page := 1
		perPage := 8

		expectedTotalItems := int64(len(expectedOrders))
		orderRepo.On("GetOrderByOrderStatus", orderStatus, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByByOrderStatus", orderStatus).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetOrderByOrderStatus(orderStatus, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Order Count by Order Status", func(t *testing.T) {
		orderStatus := "Pending"
		page := 1
		perPage := 8

		expectedErr := errors.New("gagal mendapatkan total pesanan")

		orderRepo.On("GetOrderByOrderStatus", orderStatus, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByByOrderStatus", orderStatus).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByOrderStatus(orderStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {
		orderStatus := "Pending"
		page := 1
		perPage := 8

		expectedErr := errors.New("pesanan tidak ditemukan")

		orderRepo.On("GetOrderByOrderStatus", orderStatus, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByOrderStatus(orderStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrderByDateRangeAndStatus(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	orderStatus := "Pending"
	filterType := "Minggu Ini"
	page := 1
	perPage := 8

	startOfWeek := time.Now().In(time.UTC).AddDate(0, 0, -int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)
	endOfWeek := time.Now().In(time.UTC).AddDate(0, 0, 7-int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)

	expectedOrders := []*entities.OrderModels{
		{ID: "order123", OrderStatus: "Pending"},
		{ID: "order456", OrderStatus: "Shipped"},
	}

	t.Run("Success Case - Filter by Week and Status", func(t *testing.T) {

		expectedTotalItems := int64(len(expectedOrders))
		orderRepo.On("GetOrderByDateRangeAndStatus", startOfWeek, endOfWeek, orderStatus, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRangeAndStatus", startOfWeek, endOfWeek, orderStatus).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndStatus(filterType, orderStatus, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Invalid Filter Type", func(t *testing.T) {
		invalidFilter := "InvalidFilter"

		result, totalItems, err := orderService.GetOrderByDateRangeAndStatus(invalidFilter, orderStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.EqualError(t, err, "tipe filter tidak valid")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pesanan tidak ditemukan")
		orderRepo.On("GetOrderByDateRangeAndStatus", startOfWeek, endOfWeek, orderStatus, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndStatus(filterType, orderStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case -  Error Getting Total Order Count by Date Range And Order Status", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan total pesanan")
		orderRepo.On("GetOrderByDateRangeAndStatus", startOfWeek, endOfWeek, orderStatus, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRangeAndStatus", startOfWeek, endOfWeek, orderStatus).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndStatus(filterType, orderStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrderByDateRangeAndStatusAndSearch(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	orderStatus := "Pending"
	search := "name"
	filterType := "Minggu Ini"
	page := 1
	perPage := 8

	startOfWeek := time.Now().In(time.UTC).AddDate(0, 0, -int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)
	endOfWeek := time.Now().In(time.UTC).AddDate(0, 0, 7-int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)

	expectedOrders := []*entities.OrderModels{
		{ID: "order123", OrderStatus: "Pending"},
		{ID: "order456", OrderStatus: "Shipped"},
	}

	t.Run("Success Case - Order Found", func(t *testing.T) {

		orderRepo.On("GetOrderByDateRangeAndStatusAndSearch", startOfWeek, endOfWeek, orderStatus, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRangeAndStatusAndSearch", startOfWeek, endOfWeek, orderStatus, search).Return(int64(len(expectedOrders)), nil).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndStatusAndSearch(filterType, orderStatus, search, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, int64(len(expectedOrders)), totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Invalid Filter Type", func(t *testing.T) {
		invalidFilter := "InvalidFilter"

		result, totalItems, err := orderService.GetOrderByDateRangeAndStatusAndSearch(invalidFilter, orderStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.EqualError(t, err, "tipe filter tidak valid")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pesanan tidak ditemukan")
		orderRepo.On("GetOrderByDateRangeAndStatusAndSearch", startOfWeek, endOfWeek, orderStatus, search, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndStatusAndSearch(filterType, orderStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case -  Error Getting Total Order Count by Date Range And Order Status And Search", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan total pesanan")
		orderRepo.On("GetOrderByDateRangeAndStatusAndSearch", startOfWeek, endOfWeek, orderStatus, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRangeAndStatusAndSearch", startOfWeek, endOfWeek, orderStatus, search).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndStatusAndSearch(filterType, orderStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

}

func TestOrderService_GetOrderBySearchAndDateRange(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	search := "name"
	filterType := "Minggu Ini"
	page := 1
	perPage := 8

	startOfWeek := time.Now().In(time.UTC).AddDate(0, 0, -int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)
	endOfWeek := time.Now().In(time.UTC).AddDate(0, 0, 7-int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)

	expectedOrders := []*entities.OrderModels{
		{ID: "order123", OrderStatus: "Pending"},
		{ID: "order456", OrderStatus: "Shipped"},
	}

	t.Run("Success Case - Order Found", func(t *testing.T) {

		orderRepo.On("GetOrdersBySearchAndDateRange", startOfWeek, endOfWeek, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountBySearchAndDateRange", startOfWeek, endOfWeek, search).Return(int64(len(expectedOrders)), nil).Once()

		result, totalItems, err := orderService.GetOrderBySearchAndDateRange(filterType, search, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, int64(len(expectedOrders)), totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Invalid Filter Type", func(t *testing.T) {
		invalidFilter := "InvalidFilter"

		result, totalItems, err := orderService.GetOrderBySearchAndDateRange(invalidFilter, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.EqualError(t, err, "tipe filter tidak valid")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pesanan tidak ditemukan")
		orderRepo.On("GetOrdersBySearchAndDateRange", startOfWeek, endOfWeek, search, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrderBySearchAndDateRange(filterType, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case -  Error Getting Total Order Count by Date Range And Search", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan total pesanan")
		orderRepo.On("GetOrdersBySearchAndDateRange", startOfWeek, endOfWeek, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountBySearchAndDateRange", startOfWeek, endOfWeek, search).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrderBySearchAndDateRange(filterType, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrdersBySearchAndStatus(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	search := "name"
	orderStatus := "Pending"
	page := 1
	perPage := 8

	expectedOrders := []*entities.OrderModels{
		{ID: "order123", OrderStatus: "Pending"},
		{ID: "order456", OrderStatus: "Shipped"},
	}

	t.Run("Success Case - Order Found", func(t *testing.T) {

		orderRepo.On("GetOrdersBySearchAndStatus", orderStatus, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrdersCountBySearchAndStatus", orderStatus, search).Return(int64(len(expectedOrders)), nil).Once()

		result, totalItems, err := orderService.GetOrdersBySearchAndStatus(orderStatus, search, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, int64(len(expectedOrders)), totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pesanan tidak ditemukan")
		orderRepo.On("GetOrdersBySearchAndStatus", orderStatus, search, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrdersBySearchAndStatus(orderStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case -  Error Getting Total Order Count by Date Range And Search", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan total pesanan")
		orderRepo.On("GetOrdersBySearchAndStatus", orderStatus, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrdersCountBySearchAndStatus", orderStatus, search).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrdersBySearchAndStatus(orderStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrderByPaymentStatus(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	paymentStatus := "Confirm"
	page := 1
	perPage := 8

	expectedOrders := []*entities.OrderModels{
		{ID: "order123", PaymentStatus: "Pending"},
		{ID: "order456", PaymentStatus: "Confirm"},
	}

	t.Run("Success Case - Get Orders by Payment Status", func(t *testing.T) {
		expectedTotalItems := int64(len(expectedOrders))
		orderRepo.On("GetOrderByPaymentStatus", paymentStatus, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByByPaymentStatus", paymentStatus).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetOrderByPaymentStatus(paymentStatus, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Order Count by  Payment Status", func(t *testing.T) {

		expectedErr := errors.New("gagal mendapatkan total pembayaran")

		orderRepo.On("GetOrderByPaymentStatus", paymentStatus, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByByPaymentStatus", paymentStatus).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByPaymentStatus(paymentStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pembayaran tidak ditemukan")

		orderRepo.On("GetOrderByPaymentStatus", paymentStatus, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByPaymentStatus(paymentStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrderByDateRangeAndPaymentStatus(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	paymentStatus := "Confirm"
	filterType := "Minggu Ini"
	page := 1
	perPage := 8

	startOfWeek := time.Now().In(time.UTC).AddDate(0, 0, -int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)
	endOfWeek := time.Now().In(time.UTC).AddDate(0, 0, 7-int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)

	expectedOrders := []*entities.OrderModels{
		{ID: "order123", PaymentStatus: "Pending"},
		{ID: "order456", PaymentStatus: "Confirm"},
	}

	t.Run("Success Case - Filter by Week and Status", func(t *testing.T) {

		expectedTotalItems := int64(len(expectedOrders))
		orderRepo.On("GetOrderByDateRangeAndPaymentStatus", startOfWeek, endOfWeek, paymentStatus, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRangeAndPaymentStatus", startOfWeek, endOfWeek, paymentStatus).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndPaymentStatus(filterType, paymentStatus, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Invalid Filter Type", func(t *testing.T) {
		invalidFilter := "InvalidFilter"

		result, totalItems, err := orderService.GetOrderByDateRangeAndPaymentStatus(invalidFilter, paymentStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.EqualError(t, err, "tipe filter tidak valid")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pembayaran tidak ditemukan")
		orderRepo.On("GetOrderByDateRangeAndPaymentStatus", startOfWeek, endOfWeek, paymentStatus, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndPaymentStatus(filterType, paymentStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case -  Error Getting Total Order Count by Date Range And Payment Status", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan total pembayaran")
		orderRepo.On("GetOrderByDateRangeAndPaymentStatus", startOfWeek, endOfWeek, paymentStatus, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRangeAndPaymentStatus", startOfWeek, endOfWeek, paymentStatus).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndPaymentStatus(filterType, paymentStatus, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrderByDateRangeAndPaymentStatusAndSearch(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	paymentStatus := "Confirm"
	search := "name"
	filterType := "Minggu Ini"
	page := 1
	perPage := 8

	startOfWeek := time.Now().In(time.UTC).AddDate(0, 0, -int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)
	endOfWeek := time.Now().In(time.UTC).AddDate(0, 0, 7-int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour)

	expectedOrders := []*entities.OrderModels{
		{ID: "order123", PaymentStatus: "Pending"},
		{ID: "order456", PaymentStatus: "Confirm"},
	}

	t.Run("Success Case - Order Found", func(t *testing.T) {

		orderRepo.On("GetOrderByDateRangeAndPaymentStatusAndSearch", startOfWeek, endOfWeek, paymentStatus, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRangeAndPaymentStatusAndSearch", startOfWeek, endOfWeek, paymentStatus, search).Return(int64(len(expectedOrders)), nil).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndPaymentStatusAndSearch(filterType, paymentStatus, search, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, int64(len(expectedOrders)), totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Invalid Filter Type", func(t *testing.T) {
		invalidFilter := "InvalidFilter"

		result, totalItems, err := orderService.GetOrderByDateRangeAndPaymentStatusAndSearch(invalidFilter, paymentStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.EqualError(t, err, "tipe filter tidak valid")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pembayaran tidak ditemukan")
		orderRepo.On("GetOrderByDateRangeAndPaymentStatusAndSearch", startOfWeek, endOfWeek, paymentStatus, search, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndPaymentStatusAndSearch(filterType, paymentStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failed Case -  Error Getting Total Order Count by Date Range And Order Status And Search", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan total pembayaran")
		orderRepo.On("GetOrderByDateRangeAndPaymentStatusAndSearch", startOfWeek, endOfWeek, paymentStatus, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountByDateRangeAndPaymentStatusAndSearch", startOfWeek, endOfWeek, paymentStatus, search).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrderByDateRangeAndPaymentStatusAndSearch(filterType, paymentStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

}

func TestOrderService_GetOrdersBySearchAndPaymentStatus(t *testing.T) {
	orderService, orderRepo, _, _, _, _, _, _, _, _ := setupOrderService(t)
	paymentStatus := "Confirm"
	search := "name"
	page := 1
	perPage := 8

	expectedOrders := []*entities.OrderModels{
		{ID: "order123", PaymentStatus: "Pending"},
		{ID: "order456", PaymentStatus: "Confirm"},
	}

	t.Run("Success Case - Order Found", func(t *testing.T) {

		orderRepo.On("GetOrderBySearchAndPaymentStatus", paymentStatus, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountBySearchAndPaymentStatus", paymentStatus, search).Return(int64(len(expectedOrders)), nil).Once()

		result, totalItems, err := orderService.GetOrdersBySearchAndPaymentStatus(paymentStatus, search, page, perPage)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(expectedOrders), len(result))
		assert.Equal(t, int64(len(expectedOrders)), totalItems)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pembayaran tidak ditemukan")
		orderRepo.On("GetOrderBySearchAndPaymentStatus", paymentStatus, search, 0, perPage).Return(nil, expectedErr).Once()

		result, totalItems, err := orderService.GetOrdersBySearchAndPaymentStatus(paymentStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Failure Case -  Error Getting Total Order Count by Date Range And Search", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan total pembayaran")
		orderRepo.On("GetOrderBySearchAndPaymentStatus", paymentStatus, search, 0, perPage).Return(expectedOrders, nil).Once()
		orderRepo.On("GetOrderCountBySearchAndPaymentStatus", paymentStatus, search).Return(int64(0), expectedErr).Once()

		result, totalItems, err := orderService.GetOrdersBySearchAndPaymentStatus(paymentStatus, search, page, perPage)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		orderRepo.AssertExpectations(t)
	})
}

func TestOrderService_SendNotificationPayment(t *testing.T) {
	mockUser := &entities.UserModels{
		ID:          1,
		Name:        "John",
		DeviceToken: "device_token_1",
	}

	mockOrder := &entities.OrderModels{
		IdOrder: "order_id_1",
	}

	mockFcm := &entities.FcmModels{
		ID:        1,
		OrderID:   "order_id_1",
		UserID:    1,
		Title:     "Status Pembayaran",
		Body:      "Alloo, John! Pesananmu dengan ID order_id_1 udah berhasil dibuat, nih. Ditunggu yupp!!",
		CreatedAt: time.Now(),
	}

	orderService, orderRepo, userRepo, _, _, _, _, _, fcmRepo, _ := setupOrderService(t)

	t.Run("Success Case - SendNotificationPayment", func(t *testing.T) {
		userRepo.On("GetUsersById", uint64(1)).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", "order_id_1").Return(mockOrder, nil).Once()
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		request := dto.SendNotificationPaymentRequest{
			UserID:        1,
			OrderID:       "order_id_1",
			PaymentStatus: "Menunggu Konfirmasi",
		}

		msg, err := orderService.SendNotificationPayment(request)

		assert.NoError(t, err)
		assert.Equal(t, "Alloo, John! Pesananmu dengan ID order_id_1 udah berhasil dibuat, nih. Ditunggu yupp!!", msg)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Success Case - SendNotificationPayment - Konfirmasi", func(t *testing.T) {
		userRepo.On("GetUsersById", uint64(1)).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", "order_id_1").Return(mockOrder, nil).Once()
		mockFcm.Body = "Thengkyuu, John! Pembayaran untuk pesananmu dengan ID order_id_1 udah kami terima, nih. Semoga harimu menyenangkan!"
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		request := dto.SendNotificationPaymentRequest{
			UserID:        1,
			OrderID:       "order_id_1",
			PaymentStatus: "Konfirmasi",
		}

		msg, err := orderService.SendNotificationPayment(request)

		assert.NoError(t, err)
		assert.Equal(t, "Thengkyuu, John! Pembayaran untuk pesananmu dengan ID order_id_1 udah kami terima, nih. Semoga harimu menyenangkan!", msg)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Success Case - SendNotificationPayment - Gagal", func(t *testing.T) {
		userRepo.On("GetUsersById", uint64(1)).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", "order_id_1").Return(mockOrder, nil).Once()
		mockFcm.Body = "Maaf, John. Pembayaran untuk pesanan dengan ID order_id_1 gagal, nih. Beritahu kami apabila kamu butuh bantuan yaa!!"
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		request := dto.SendNotificationPaymentRequest{
			UserID:        1,
			OrderID:       "order_id_1",
			PaymentStatus: "Gagal",
		}

		msg, err := orderService.SendNotificationPayment(request)

		assert.NoError(t, err)
		assert.Equal(t, "Maaf, John. Pembayaran untuk pesanan dengan ID order_id_1 gagal, nih. Beritahu kami apabila kamu butuh bantuan yaa!!", msg)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - Invalid Payment Status", func(t *testing.T) {
		invalidStatus := "Invalid Status"

		userRepo.On("GetUsersById", uint64(1)).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", "order_id_1").Return(mockOrder, nil).Once()

		request := dto.SendNotificationPaymentRequest{
			UserID:        1,
			OrderID:       "order_id_1",
			PaymentStatus: invalidStatus,
		}

		_, err := orderService.SendNotificationPayment(request)

		expectedErr := errors.New("Status pesanan tidak valid")
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		notFoundUserID := uint64(999)
		orderID := "order_id_1"
		paymentStatus := "Menunggu Konfirmasi"

		expectedErr := errors.New("pengguna tidak ditemukan")
		userRepo.On("GetUsersById", notFoundUserID).Return(nil, expectedErr).Once()

		request := dto.SendNotificationPaymentRequest{
			UserID:        notFoundUserID,
			OrderID:       orderID,
			PaymentStatus: paymentStatus,
		}

		_, err := orderService.SendNotificationPayment(request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - Order Not Found", func(t *testing.T) {
		userID := uint64(999)
		orderID := "order_id_1"
		paymentStatus := "Menunggu Konfirmasi"

		expectedErr := errors.New("pesanan tidak ditemukan")
		userRepo.On("GetUsersById", userID).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", "order_id_1").Return(nil, expectedErr).Once()

		request := dto.SendNotificationPaymentRequest{
			UserID:        userID,
			OrderID:       orderID,
			PaymentStatus: paymentStatus,
		}

		_, err := orderService.SendNotificationPayment(request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - Failed Create FCM", func(t *testing.T) {
		expectedErr := errors.New("failed to create FCM")

		userRepo.On("GetUsersById", uint64(1)).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", "order_id_1").Return(mockOrder, nil).Once()
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(nil, expectedErr).Once()

		request := dto.SendNotificationPaymentRequest{
			UserID:        1,
			OrderID:       "order_id_1",
			PaymentStatus: "Konfirmasi",
		}

		_, err := orderService.SendNotificationPayment(request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})
}

func TestOrderService_SendNotificationOrder(t *testing.T) {
	userID := uint64(999)
	orderID := "order_id_1"
	orderStatus := "Menunggu Konfirmasi"

	mockUser := &entities.UserModels{
		ID:          1,
		Name:        "John",
		DeviceToken: "device_token_1",
	}

	mockOrder := &entities.OrderModels{
		IdOrder: "order_id_1",
	}

	mockFcm := &entities.FcmModels{
		ID:        1,
		OrderID:   "order_id_1",
		UserID:    1,
		Title:     "Status Pengiriman",
		Body:      fmt.Sprintf("Alloo, %s! Pesananmu dengan ID %s sedang menunggu konfirmasi, nih. Ditunggu yupp!", mockUser.Name, mockOrder.IdOrder),
		CreatedAt: time.Now(),
	}

	orderService, orderRepo, userRepo, _, _, _, _, _, fcmRepo, _ := setupOrderService(t)

	t.Run("Success Case - SendNotificationOrder", func(t *testing.T) {
		userRepo.On("GetUsersById", uint64(1)).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", "order_id_1").Return(mockOrder, nil).Once()
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      1,
			OrderID:     "order_id_1",
			OrderStatus: "Menunggu Konfirmasi",
		}

		msg, err := orderService.SendNotificationOrder(request)

		assert.NoError(t, err)
		assert.Equal(t, "Alloo, John! Pesananmu dengan ID order_id_1 sedang menunggu konfirmasi, nih. Ditunggu yupp!", msg)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - Order Not Found", func(t *testing.T) {

		expectedErr := errors.New("pesanan tidak ditemukan")
		userRepo.On("GetUsersById", userID).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", "order_id_1").Return(nil, expectedErr).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      userID,
			OrderID:     orderID,
			OrderStatus: orderStatus,
		}

		_, err := orderService.SendNotificationOrder(request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		notFoundUserID := uint64(999)
		expectedErr := errors.New("pengguna tidak ditemukan")
		userRepo.On("GetUsersById", notFoundUserID).Return(nil, expectedErr).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      notFoundUserID,
			OrderID:     orderID,
			OrderStatus: orderStatus,
		}

		_, err := orderService.SendNotificationOrder(request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - Failed Create FCM", func(t *testing.T) {
		expectedErr := errors.New("failed to create FCM")

		userRepo.On("GetUsersById", userID).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", orderID).Return(mockOrder, nil).Once()
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(nil, expectedErr).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      userID,
			OrderID:     orderID,
			OrderStatus: orderStatus,
		}

		_, err := orderService.SendNotificationOrder(request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - Invalid Order Status", func(t *testing.T) {
		invalidStatus := "Invalid Status"

		userRepo.On("GetUsersById", userID).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", orderID).Return(mockOrder, nil).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      userID,
			OrderID:     orderID,
			OrderStatus: invalidStatus,
		}

		_, err := orderService.SendNotificationOrder(request)

		expectedErr := errors.New("Status pengiriman tidak valid")
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Success Case - SendNotificationOrder - Gagal", func(t *testing.T) {
		userRepo.On("GetUsersById", userID).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", orderID).Return(mockOrder, nil).Once()
		mockFcm.Body = fmt.Sprintf("Sowwy, %s. Pesananmu dengan ID %s gagal. Coba lagi, yukk!", mockUser.Name, mockOrder.IdOrder)
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      userID,
			OrderID:     orderID,
			OrderStatus: "Gagal",
		}

		msg, err := orderService.SendNotificationOrder(request)

		assert.NoError(t, err)
		assert.Equal(t, "Sowwy, John. Pesananmu dengan ID order_id_1 gagal. Coba lagi, yukk!", msg)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Success Case - SendNotificationOrder - Proses", func(t *testing.T) {
		userRepo.On("GetUsersById", userID).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", orderID).Return(mockOrder, nil).Once()
		mockFcm.Body = fmt.Sprintf("Alloo, %s! Pesananmu dengan ID %s sedang dalam proses, nih. Ditunggu yupp!", mockUser.Name, mockOrder.IdOrder)
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      userID,
			OrderID:     orderID,
			OrderStatus: "Proses",
		}

		msg, err := orderService.SendNotificationOrder(request)

		assert.NoError(t, err)
		assert.Equal(t, "Alloo, John! Pesananmu dengan ID order_id_1 sedang dalam proses, nih. Ditunggu yupp!", msg)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Success Case - SendNotificationOrder - Pengiriman", func(t *testing.T) {
		userRepo.On("GetUsersById", userID).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", orderID).Return(mockOrder, nil).Once()
		mockFcm.Body = fmt.Sprintf("Alloo, %s! Pesanan dengan ID %s udah dalam proses pengiriman, nih. Mohon ditunggu yupp!", mockUser.Name, mockOrder.IdOrder)
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      userID,
			OrderID:     orderID,
			OrderStatus: "Pengiriman",
		}

		msg, err := orderService.SendNotificationOrder(request)

		assert.NoError(t, err)
		assert.Equal(t, "Alloo, John! Pesanan dengan ID order_id_1 udah dalam proses pengiriman, nih. Mohon ditunggu yupp!", msg)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

	t.Run("Success Case - SendNotificationOrder - Selesai", func(t *testing.T) {
		userRepo.On("GetUsersById", userID).Return(mockUser, nil).Once()
		orderRepo.On("GetOrderById", orderID).Return(mockOrder, nil).Once()
		mockFcm.Body = fmt.Sprintf("Yeayy, %s! Pesananmu dengan ID %s udah sampai tujuan, nih. Semoga sukakk yupp!", mockUser.Name, mockOrder.IdOrder)
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		request := dto.SendNotificationOrderRequest{
			UserID:      userID,
			OrderID:     orderID,
			OrderStatus: "Selesai",
		}

		msg, err := orderService.SendNotificationOrder(request)

		assert.NoError(t, err)
		assert.Equal(t, "Yeayy, John! Pesananmu dengan ID order_id_1 udah sampai tujuan, nih. Semoga sukakk yupp!", msg)

		userRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)
	})

}

func TestOrderService_AcceptOrder(t *testing.T) {
	orderID := "order_id_1"
	mockOrder := &entities.OrderModels{
		ID:     "order_id_1",
		UserID: 1,
	}

	// mockUser := &entities.UserModels{
	// 	ID:          1,
	// 	DeviceToken: "user_device_token",
	// }

	orderService, orderRepo, userRepo, _, _, _, _, _, _, _ := setupOrderService(t)

	// t.Run("Success Case - Order Accepted", func(t *testing.T) {
	// 	orderRepo.On("GetOrderById", mock.AnythingOfType("string")).Return(mockOrder, nil)
	// 	userRepo.On("GetUsersById", mock.AnythingOfType("uint64")).Return(mockUser, nil)
	// 	orderRepo.On("AcceptOrder", mockOrder.ID, "Selesai").Return(nil).Once()
	// 	fcmRepo.On("SendNotificationOrder", mock.AnythingOfType("dto.SendNotificationOrderRequest")).Return("", nil).Once()

	// 	err := orderService.AcceptOrder(orderID)

	// 	assert.NoError(t, err)

	// 	orderRepo.AssertExpectations(t)
	// 	userRepo.AssertExpectations(t)
	// 	fcmRepo.AssertExpectations(t)
	// })

	t.Run("Failure Case - Order Not Found", func(t *testing.T) {
		expectedErr := errors.New("pesanan tidak ditemukan")
		orderRepo.On("GetOrderById", orderID).Return(nil, expectedErr).Once()

		err := orderService.AcceptOrder(orderID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		orderRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		orderRepo.On("GetOrderById", orderID).Return(mockOrder, nil).Once()
		userRepo.On("GetUsersById", uint64(1)).Return(nil, expectedErr).Once()

		err := orderService.AcceptOrder(orderID)

		assert.Error(t, err)
		assert.EqualError(t, err, "pengguna tidak ditemukan")

		orderRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

}

func TestOrderService_CreateOrder(t *testing.T) {
	orderService, orderRepo, userRepo, productRepo, _, voucherRepo, addressRepo, cartRepo, fcmRepo, generatorRepo :=
		setupOrderService(t)

	orderService.repo = orderRepo
	orderService.generatorID = generatorRepo
	userID := uint64(1)
	orderID := "fake_order_id"

	createOrderRequest := &dto.CreateOrderRequest{
		AddressID:     1,
		VoucherID:     1,
		Note:          "test order",
		ProductID:     1,
		Quantity:      1,
		PaymentMethod: "whatsapp",
	}

	mockAddress := &entities.AddressModels{
		ID: userID,
	}

	mockVoucher := &entities.VoucherModels{
		ID:    1,
		Stock: 10,
	}

	mockUser := &entities.UserModels{
		ID:          1,
		Name:        "John",
		DeviceToken: "device_token_1",
	}

	mockProduct := &entities.ProductModels{
		ID:    1,
		Stock: 100,
	}

	mockOrder := &entities.OrderModels{
		IdOrder: "fake_order_id",
	}

	mockFcm := &entities.FcmModels{
		ID:        1,
		OrderID:   "order_id_1",
		UserID:    1,
		Title:     "Status Pengiriman",
		Body:      fmt.Sprintf("Alloo, %s! Pesananmu dengan ID %s sedang menunggu konfirmasi, nih. Ditunggu yupp!", mockUser.Name, mockOrder.IdOrder),
		CreatedAt: time.Now(),
	}

	t.Run("Success Case - Create Order", func(t *testing.T) {
		generatorRepo.On("GenerateUUID").Return(orderID, nil)
		generatorRepo.On("GenerateOrderID").Return("fake_id_order", nil)
		addressRepo.On("GetAddressByID", createOrderRequest.AddressID).Return(mockAddress, nil)
		voucherRepo.On("GetVoucherById", createOrderRequest.VoucherID).Return(mockVoucher, nil)
		productRepo.On("GetProductByID", createOrderRequest.ProductID).Return(mockProduct, nil)
		cartRepo.On("IsProductInCart", userID, mockProduct.ID).Return(false).Once()
		productRepo.On("ReduceStockWhenPurchasing", createOrderRequest.ProductID, createOrderRequest.Quantity).Return(nil)
		orderRepo.On("CreateOrder", mock.AnythingOfType("*entities.OrderModels")).Return(mockOrder, nil)
		voucherRepo.On("DeleteUserVoucherClaims", mock.AnythingOfType("uint64"), mock.AnythingOfType("uint64")).Return(nil)
		userRepo.On("GetUsersById", mock.AnythingOfType("uint64")).Return(mockUser, nil)
		orderRepo.On("GetOrderById", mock.AnythingOfType("string")).Return(mockOrder, nil)
		fcmRepo.On("SendMessageNotification", mock.Anything).Return("", nil).Once()
		fcmRepo.On("CreateFcm", mock.Anything).Return(mockFcm, nil).Once()

		result, err := orderService.CreateOrder(userID, createOrderRequest)
		generatorRepo.AssertExpectations(t)
		addressRepo.AssertExpectations(t)
		voucherRepo.AssertExpectations(t)
		cartRepo.AssertExpectations(t)
		productRepo.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
		voucherRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		fcmRepo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Failed Case- Create order id", func(t *testing.T) {
		orderService, _, _, _, _, _, _, _, _, generatorRepo := setupOrderService(t)

		generatorRepo.On("GenerateUUID").Return("", errors.New("gagal membuat id pesanan"))

		_, err := orderService.CreateOrder(1, &dto.CreateOrderRequest{})

		generatorRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Equal(t, "gagal membuat id pesanan", err.Error())
	})

	t.Run("Failed Case - Create id_order", func(t *testing.T) {
		orderService, _, _, _, _, _, _, _, _, generatorRepo := setupOrderService(t)

		generatorRepo.On("GenerateUUID").Return("fake_order_id", nil)
		generatorRepo.On("GenerateOrderID").Return("", errors.New("gagal membuat id_order"))

		_, err := orderService.CreateOrder(1, &dto.CreateOrderRequest{})

		generatorRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Equal(t, "gagal membuat id_order", err.Error())
	})

	t.Run("Failed Case - AddressNotFound", func(t *testing.T) {
		orderService, _, _, _, _, _, addressRepo, _, _, generatorRepo := setupOrderService(t)

		generatorRepo.On("GenerateUUID").Return("fake_order_id", nil)
		generatorRepo.On("GenerateOrderID").Return("fake_id_order", nil)
		addressRepo.On("GetAddressByID", mock.AnythingOfType("uint64")).Return(nil, errors.New("alamat tidak ditemukan"))

		_, err := orderService.CreateOrder(1, &dto.CreateOrderRequest{
			AddressID: 1,
		})

		generatorRepo.AssertExpectations(t)
		addressRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Equal(t, "alamat tidak ditemukan", err.Error())
	})

	t.Run("Failed Case - VoucherNotFound", func(t *testing.T) {
		orderService, _, _, _, _, voucherRepo, addressRepo, _, _, generatorRepo := setupOrderService(t)

		generatorRepo.On("GenerateUUID").Return("fake_order_id", nil)
		generatorRepo.On("GenerateOrderID").Return("fake_id_order", nil)
		addressRepo.On("GetAddressByID", mock.AnythingOfType("uint64")).Return(&entities.AddressModels{}, nil)
		voucherRepo.On("GetVoucherById", mock.AnythingOfType("uint64")).Return(nil, errors.New("kupon tidak ditemukan"))

		_, err := orderService.CreateOrder(1, &dto.CreateOrderRequest{
			AddressID: 1,
			VoucherID: 1,
		})

		generatorRepo.AssertExpectations(t)
		addressRepo.AssertExpectations(t)
		voucherRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Equal(t, "kupon tidak ditemukan", err.Error())
	})

	t.Run("InvalidPaymentMethod", func(t *testing.T) {
		orderService, _, _, _, _, _, addressRepo, _, _, generatorRepo := setupOrderService(t)

		generatorRepo.On("GenerateUUID").Return("fake_order_id", nil)
		generatorRepo.On("GenerateOrderID").Return("fake_id_order", nil)
		addressRepo.On("GetAddressByID", mock.AnythingOfType("uint64")).Return(&entities.AddressModels{}, nil)

		_, err := orderService.CreateOrder(1, &dto.CreateOrderRequest{
			AddressID:     1,
			PaymentMethod: "invalid_method",
		})

		assert.NotNil(t, err)
		assert.Equal(t, "jenis pembayaran tidak valid", err.Error())
	})
}
