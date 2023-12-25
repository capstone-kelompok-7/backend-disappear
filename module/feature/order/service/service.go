package service

import (
	"errors"
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/sendnotif"
	"github.com/sirupsen/logrus"
	"math"
	"strings"
	"time"
)

type OrderService struct {
	repo           order.RepositoryOrderInterface
	generatorID    utils.GeneratorInterface
	productService product.ServiceProductInterface
	voucherService voucher.ServiceVoucherInterface
	addressService address.ServiceAddressInterface
	userService    users.ServiceUserInterface
	cartService    cart.ServiceCartInterface
	fcmService     fcm.ServiceFcmInterface
}

func NewOrderService(
	repo order.RepositoryOrderInterface,
	generatorID utils.GeneratorInterface,
	productService product.ServiceProductInterface,
	voucherService voucher.ServiceVoucherInterface,
	addressService address.ServiceAddressInterface,
	userService users.ServiceUserInterface,
	cartService cart.ServiceCartInterface,
	fcmService fcm.ServiceFcmInterface,
) order.ServiceOrderInterface {
	return &OrderService{
		repo:           repo,
		generatorID:    generatorID,
		productService: productService,
		voucherService: voucherService,
		addressService: addressService,
		userService:    userService,
		cartService:    cartService,
		fcmService:     fcmService,
	}
}

func (s *OrderService) GetAll(page, perPage int) ([]*entities.OrderModels, int64, error) {
	orders, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalOrderCount()
	if err != nil {
		return nil, 0, err
	}

	return orders, totalItems, nil
}

func (s *OrderService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	pageInt := page
	if pageInt <= 0 {
		pageInt = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	if pageInt > totalPages {
		pageInt = totalPages
	}

	return pageInt, totalPages
}

func (s *OrderService) GetNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}
	return totalPages
}

func (s *OrderService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}
	return 1
}

func (s *OrderService) GetOrdersByName(page, perPage int, name string) ([]*entities.OrderModels, int64, error) {
	orders, err := s.repo.FindByName(page, perPage, name)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalCustomerCountByName(name)
	if err != nil {
		return nil, 0, err
	}

	return orders, totalItems, nil
}

func (s *OrderService) GetOrderById(orderID string) (*entities.OrderModels, error) {
	orders, err := s.repo.GetOrderById(orderID)
	if err != nil {
		return nil, errors.New("gagal mendapatkan pesanan")
	}
	return orders, nil
}

func (s *OrderService) CreateOrder(userID uint64, request *dto.CreateOrderRequest) (interface{}, error) {
	orderID, err := s.generatorID.GenerateUUID()
	if err != nil {
		return nil, errors.New("gagal membuat id pesanan")
	}

	idOrder, err := s.generatorID.GenerateOrderID()
	if err != nil {
		return nil, errors.New("gagal membuat id_order")
	}

	addresses, err := s.addressService.GetAddressByID(request.AddressID)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}

	var vouchers *entities.VoucherModels
	if request.VoucherID != 0 {
		vouchers, err = s.voucherService.GetVoucherById(request.VoucherID)
		if err != nil {
			return nil, errors.New("kupon tidak ditemukan")
		}
	}

	var validPaymentMethods = map[string]bool{
		"whatsapp":      true,
		"telegram":      true,
		"qris":          true,
		"bank_transfer": true,
		"gopay":         true,
		"shopepay":      true,
	}

	if !validPaymentMethods[request.PaymentMethod] {
		return nil, errors.New("jenis pembayaran tidak valid")
	}

	var orderDetails []entities.OrderDetailsModels
	var totalQuantity, totalGramPlastic, totalExp, totalPrice, totalDiscount uint64

	products, err := s.productService.GetProductByID(request.ProductID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	if products.Stock < request.Quantity {
		return nil, errors.New("stok tidak mencukupi untuk pesanan ini")
	}

	orderDetail := entities.OrderDetailsModels{
		OrderID:          orderID,
		ProductID:        request.ProductID,
		Quantity:         request.Quantity,
		TotalGramPlastic: products.GramPlastic * request.Quantity,
		TotalExp:         products.Exp * request.Quantity,
		TotalPrice:       request.Quantity * (products.Price - products.Discount),
		TotalDiscount:    products.Discount * request.Quantity,
	}

	totalQuantity += request.Quantity
	totalGramPlastic += orderDetail.TotalGramPlastic
	totalExp += orderDetail.TotalExp
	totalPrice += orderDetail.TotalPrice
	totalDiscount += orderDetail.TotalDiscount

	orderDetails = append(orderDetails, orderDetail)

	if isInCart := s.cartService.IsProductInCart(userID, products.ID); isInCart {
		if err := s.cartService.RemoveProductFromCart(userID, products.ID); err != nil {
			return nil, errors.New("gagal menghapus keranjang")
		}
	}

	if err := s.productService.ReduceStockWhenPurchasing(request.ProductID, request.Quantity); err != nil {
		return nil, errors.New("gagal mengurangi stok")
	}

	var discountFromVoucher uint64
	if request.VoucherID != 0 && totalPrice >= vouchers.MinPurchase {
		discountFromVoucher = vouchers.Discount
	}

	var voucherID *uint64
	if request.VoucherID != 0 {
		voucherID = &request.VoucherID
	}

	grandTotalPrice := totalPrice
	totalAmountPaid := grandTotalPrice + 2000 + 0 - discountFromVoucher

	newData := &entities.OrderModels{
		ID:                    orderID,
		IdOrder:               idOrder,
		AddressID:             addresses.ID,
		UserID:                userID,
		VoucherID:             voucherID,
		Note:                  request.Note,
		GrandTotalGramPlastic: totalGramPlastic,
		GrandTotalExp:         totalExp,
		GrandTotalQuantity:    totalQuantity,
		GrandTotalPrice:       grandTotalPrice,
		ShipmentFee:           0,
		AdminFees:             2000,
		GrandTotalDiscount:    totalDiscount,
		TotalAmountPaid:       totalAmountPaid,
		OrderStatus:           "Menunggu Konfirmasi",
		PaymentStatus:         "Menunggu Konfirmasi",
		PaymentMethod:         request.PaymentMethod,
		StatusOrderDate:       time.Now(),
		CreatedAt:             time.Now(),
		OrderDetails:          orderDetails,
	}
	createdOrder, err := s.repo.CreateOrder(newData)
	if err != nil {
		return nil, err
	}

	if request.VoucherID != 0 {
		if err := s.voucherService.DeleteVoucherClaims(userID, vouchers.ID); err != nil {
			return nil, err
		}
	}

	user, err := s.userService.GetUsersById(createdOrder.UserID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	notificationRequest := dto.SendNotificationPaymentRequest{
		OrderID:       createdOrder.ID,
		UserID:        createdOrder.UserID,
		PaymentStatus: "Menunggu Konfirmasi",
		Token:         user.DeviceToken,
	}
	_, err = s.SendNotificationPayment(notificationRequest)
	if err != nil {
		logrus.Error("Gagal mengirim notifikasi: ", err)
		return nil, err
	}

	switch request.PaymentMethod {
	case "whatsapp", "telegram":
		return s.ProcessManualPayment(orderID)
	case "qris", "bank_transfer", "gopay", "shopepay":
		return s.ProcessGatewayPayment(totalAmountPaid, createdOrder.ID, request.PaymentMethod, user.Name, user.Email)
	default:
		return nil, errors.New("jenis pembayaran tidak valid")
	}
}

func (s *OrderService) CreateOrderFromCart(userID uint64, request *dto.CreateOrderCartRequest) (interface{}, error) {
	orderID, err := s.generatorID.GenerateUUID()
	if err != nil {
		return nil, errors.New("gagal membuat id pesanan")
	}

	idOrder, err := s.generatorID.GenerateOrderID()
	if err != nil {
		return nil, errors.New("gagal membuat id_order")
	}

	addresses, err := s.addressService.GetAddressByID(request.AddressID)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}

	var vouchers *entities.VoucherModels
	if request.VoucherID != 0 {
		vouchers, err = s.voucherService.GetVoucherById(request.VoucherID)
		if err != nil {
			return nil, errors.New("kupon tidak ditemukan")
		}
	}

	var validPaymentMethods = map[string]bool{
		"whatsapp":      true,
		"telegram":      true,
		"qris":          true,
		"bank_transfer": true,
		"gopay":         true,
		"shopepay":      true,
	}

	if !validPaymentMethods[request.PaymentMethod] {
		return nil, errors.New("jenis pembayaran tidak valid")
	}

	var cartItems []*entities.CartItemModels
	for _, itemID := range request.CartItems {
		cartItem, err := s.cartService.GetCartItems(itemID.ID)
		if err != nil {
			return nil, errors.New("gagal mendapatkan detail item keranjang")
		}
		cartItems = append(cartItems, cartItem)
	}

	var orderDetails []entities.OrderDetailsModels
	var totalQuantity, totalGramPlastic, totalExp, totalPrice, totalDiscount uint64

	for _, cartItem := range cartItems {
		products, err := s.productService.GetProductByID(cartItem.ProductID)
		if err != nil {
			return nil, errors.New("produk tidak ditemukan")
		}

		if products.Stock < cartItem.Quantity {
			return nil, errors.New("stok tidak mencukupi untuk pesanan ini")
		}

		orderDetail := entities.OrderDetailsModels{
			OrderID:          orderID,
			ProductID:        cartItem.ProductID,
			Quantity:         cartItem.Quantity,
			TotalGramPlastic: products.GramPlastic * cartItem.Quantity,
			TotalExp:         products.Exp * cartItem.Quantity,
			TotalPrice:       cartItem.Quantity * (products.Price - products.Discount),
			TotalDiscount:    products.Discount * cartItem.Quantity,
		}

		totalQuantity += cartItem.Quantity
		totalGramPlastic += orderDetail.TotalGramPlastic
		totalExp += orderDetail.TotalExp
		totalPrice += orderDetail.TotalPrice
		totalDiscount += orderDetail.TotalDiscount

		orderDetails = append(orderDetails, orderDetail)

		if err := s.cartService.DeleteCartItem(cartItem.ID); err != nil {
			return nil, errors.New("gagal menghapus produk dari keranjang")
		}

		if err := s.productService.ReduceStockWhenPurchasing(cartItem.ProductID, cartItem.Quantity); err != nil {
			return nil, errors.New("gagal mengurangi stok produk")
		}
	}

	var discountFromVoucher uint64
	if request.VoucherID != 0 && totalPrice >= vouchers.MinPurchase {
		discountFromVoucher = vouchers.Discount
	}

	var voucherID *uint64
	if request.VoucherID != 0 {
		voucherID = &request.VoucherID
	}

	grandTotalPrice := totalPrice
	totalAmountPaid := grandTotalPrice + 2000 + 0 - discountFromVoucher

	newData := &entities.OrderModels{
		ID:                    orderID,
		IdOrder:               idOrder,
		AddressID:             addresses.ID,
		UserID:                userID,
		VoucherID:             voucherID,
		Note:                  request.Note,
		GrandTotalGramPlastic: totalGramPlastic,
		GrandTotalExp:         totalExp,
		GrandTotalQuantity:    totalQuantity,
		GrandTotalPrice:       grandTotalPrice,
		ShipmentFee:           0,
		AdminFees:             2000,
		GrandTotalDiscount:    totalDiscount,
		TotalAmountPaid:       totalAmountPaid,
		OrderStatus:           "Menunggu Konfirmasi",
		PaymentStatus:         "Menunggu Konfirmasi",
		PaymentMethod:         request.PaymentMethod,
		StatusOrderDate:       time.Now(),
		CreatedAt:             time.Now(),
		OrderDetails:          orderDetails,
	}

	createdOrder, err := s.repo.CreateOrder(newData)
	if err != nil {
		return nil, err
	}

	if request.VoucherID != 0 {
		if err := s.voucherService.DeleteVoucherClaims(userID, vouchers.ID); err != nil {
			return nil, err
		}
	}
	user, err := s.userService.GetUsersById(createdOrder.UserID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	notificationRequest := dto.SendNotificationPaymentRequest{
		OrderID:       createdOrder.ID,
		UserID:        createdOrder.UserID,
		PaymentStatus: "Menunggu Konfirmasi",
		Token:         user.DeviceToken,
	}
	_, err = s.SendNotificationPayment(notificationRequest)
	if err != nil {
		logrus.Error("Gagal mengirim notifikasi: ", err)
		return nil, err
	}

	switch request.PaymentMethod {
	case "whatsapp", "telegram":
		return s.ProcessManualPayment(orderID)
	case "qris", "bank_transfer", "gopay", "shopepay":
		return s.ProcessGatewayPayment(totalAmountPaid, createdOrder.ID, request.PaymentMethod, user.Name, user.Email)
	default:
		return nil, errors.New("jenis pembayaran tidak valid")
	}
}

func (s *OrderService) ProcessManualPayment(orderID string) (*entities.OrderModels, error) {
	result, err := s.repo.GetOrderById(orderID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) ProcessGatewayPayment(totalAmountPaid uint64, orderID string, paymentMethod, name, email string) (interface{}, error) {
	result, err := s.repo.ProcessGatewayPayment(totalAmountPaid, orderID, paymentMethod, name, email)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) CallBack(notifPayload map[string]interface{}) error {
	orderID, exist := notifPayload["order_id"].(string)
	if !exist {
		return errors.New("invalid notification payload")
	}

	status, err := s.repo.CheckTransaction(orderID)
	if err != nil {
		return err
	}

	transaction, err := s.repo.GetOrderById(orderID)
	if err != nil {
		return errors.New("transaction data not found")
	}

	if status.PaymentStatus == "Konfirmasi" {
		if err := s.ConfirmPayment(transaction.ID); err != nil {
			return err
		}
	} else if status.PaymentStatus == "Gagal" {
		if err := s.CancelPayment(transaction.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *OrderService) ConfirmPayment(orderID string) error {
	orders, err := s.repo.GetOrderById(orderID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}

	orders.OrderStatus = "Proses"
	orders.PaymentStatus = "Konfirmasi"

	if err := s.repo.ConfirmPayment(orders.ID, orders.OrderStatus, orders.PaymentStatus); err != nil {
		return err
	}

	user, err := s.userService.GetUsersById(orders.UserID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	user.Exp += orders.GrandTotalExp
	if _, err := s.userService.UpdateUserExp(user.ID, user.Exp); err != nil {
		return err
	}

	user.TotalGram += orders.GrandTotalGramPlastic
	if _, err := s.userService.UpdateUserContribution(user.ID, user.TotalGram); err != nil {
		return err
	}

	notificationRequest := dto.SendNotificationPaymentRequest{
		OrderID:       orderID,
		UserID:        user.ID,
		PaymentStatus: "Konfirmasi",
		Token:         user.DeviceToken,
	}
	_, err = s.SendNotificationPayment(notificationRequest)
	if err != nil {
		logrus.Error("Gagal mengirim notifikasi: ", err)
		return err
	}

	return nil
}

func (s *OrderService) CancelPayment(orderID string) error {
	orders, err := s.repo.GetOrderById(orderID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}

	orders.OrderStatus = "Gagal"
	orders.PaymentStatus = "Gagal"

	for _, orderDetail := range orders.OrderDetails {
		if err := s.productService.IncreaseStock(orderDetail.ProductID, orderDetail.Quantity); err != nil {
			return errors.New("gagal menambah stok produk")
		}
	}

	if err := s.repo.ConfirmPayment(orderID, orders.OrderStatus, orders.PaymentStatus); err != nil {
		return errors.New("gagal membatalkan pesanan")
	}

	user, err := s.userService.GetUsersById(orders.UserID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	notificationRequest := dto.SendNotificationPaymentRequest{
		OrderID:       orderID,
		UserID:        user.ID,
		PaymentStatus: "Gagal",
		Token:         user.DeviceToken,
	}
	_, err = s.SendNotificationPayment(notificationRequest)
	if err != nil {
		logrus.Error("Gagal mengirim notifikasi: ", err)
		return err
	}

	return nil
}

func (s *OrderService) UpdateOrderStatus(req *dto.UpdateOrderStatus) error {
	orders, err := s.repo.GetOrderById(req.OrderID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}

	if err := s.repo.UpdateOrderStatus(req); err != nil {
		return err
	}

	user, err := s.userService.GetUsersById(orders.UserID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	notificationRequest := dto.SendNotificationOrderRequest{
		OrderID:     orders.ID,
		UserID:      user.ID,
		OrderStatus: req.OrderStatus,
		Token:       user.DeviceToken,
	}
	_, err = s.SendNotificationOrder(notificationRequest)
	if err != nil {
		logrus.Error("Gagal mengirim notifikasi: ", err)
		return err
	}

	return nil
}

func (s *OrderService) GetAllOrdersByUserID(userID uint64) ([]*entities.OrderModels, error) {
	user, err := s.userService.GetUsersById(userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	result, err := s.repo.GetAllOrdersByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) GetAllOrdersWithFilter(userID uint64, orderStatus string) ([]*entities.OrderModels, error) {
	user, err := s.userService.GetUsersById(userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	result, err := s.repo.GetAllOrdersWithFilter(user.ID, orderStatus)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) AcceptOrder(orderID string) error {
	orders, err := s.repo.GetOrderById(orderID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}

	user, err := s.userService.GetUsersById(orders.UserID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	orders.OrderStatus = "Selesai"

	if err := s.repo.AcceptOrder(orders.ID, orders.OrderStatus); err != nil {
		return err
	}

	notificationRequest := dto.SendNotificationOrderRequest{
		OrderID:     orders.ID,
		UserID:      user.ID,
		OrderStatus: "Selesai",
		Token:       user.DeviceToken,
	}
	_, err = s.SendNotificationOrder(notificationRequest)
	if err != nil {
		logrus.Error("Gagal mengirim notifikasi: ", err)
		return err
	}

	return nil
}

func (s *OrderService) Tracking(courier, awb string) (map[string]interface{}, error) {
	result, err := s.repo.Tracking(courier, awb)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) GetOrderByDateRange(filterType string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	startDate, endDate, err := s.GetFilterDateRange(filterType)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	result, err := s.repo.GetOrderByDateRange(startDate, endDate, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pesanan tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountByDateRange(startDate, endDate)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pesanan")
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrderByOrderStatus(orderStatus string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	offset := (page - 1) * perPage

	result, err := s.repo.GetOrderByOrderStatus(orderStatus, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pesanan tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountByByOrderStatus(orderStatus)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pesanan")
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrderByDateRangeAndStatus(filterType, status string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	startDate, endDate, err := s.GetFilterDateRange(filterType)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	result, err := s.repo.GetOrderByDateRangeAndStatus(startDate, endDate, status, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pesanan tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountByDateRangeAndStatus(startDate, endDate, status)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pesanan")
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrderByDateRangeAndStatusAndSearch(filterType, status, search string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	startDate, endDate, err := s.GetFilterDateRange(filterType)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	result, err := s.repo.GetOrderByDateRangeAndStatusAndSearch(startDate, endDate, status, search, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pesanan tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountByDateRangeAndStatusAndSearch(startDate, endDate, status, search)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pesanan")
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrderBySearchAndDateRange(filterType, search string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	startDate, endDate, err := s.GetFilterDateRange(filterType)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	result, err := s.repo.GetOrdersBySearchAndDateRange(startDate, endDate, search, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pesanan tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountBySearchAndDateRange(startDate, endDate, search)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pesanan")
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrdersBySearchAndStatus(status, search string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	offset := (page - 1) * perPage

	result, err := s.repo.GetOrdersBySearchAndStatus(status, search, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pesanan tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrdersCountBySearchAndStatus(status, search)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pesanan")
	}

	return result, totalItems, nil
}

// Payment
func (s *OrderService) GetOrderByPaymentStatus(orderStatus string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	offset := (page - 1) * perPage

	result, err := s.repo.GetOrderByPaymentStatus(orderStatus, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pembayaran tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountByByPaymentStatus(orderStatus)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pembayaran")
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrderByDateRangeAndPaymentStatus(filterType, status string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	startDate, endDate, err := s.GetFilterDateRange(filterType)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	result, err := s.repo.GetOrderByDateRangeAndPaymentStatus(startDate, endDate, status, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pembayaran tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountByDateRangeAndPaymentStatus(startDate, endDate, status)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pembayaran")
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrderByDateRangeAndPaymentStatusAndSearch(filterType, status, search string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	startDate, endDate, err := s.GetFilterDateRange(filterType)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	result, err := s.repo.GetOrderByDateRangeAndPaymentStatusAndSearch(startDate, endDate, status, search, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pembayaran tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountByDateRangeAndPaymentStatusAndSearch(startDate, endDate, status, search)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pembayaran")
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrdersBySearchAndPaymentStatus(status, search string, page, perPage int) ([]*entities.OrderModels, int64, error) {
	offset := (page - 1) * perPage

	result, err := s.repo.GetOrderBySearchAndPaymentStatus(status, search, offset, perPage)
	if err != nil {
		return nil, 0, errors.New("pembayaran tidak ditemukan")
	}

	totalItems, err := s.repo.GetOrderCountBySearchAndPaymentStatus(status, search)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total pembayaran")
	}

	return result, totalItems, nil
}

func (s *OrderService) SendNotificationPayment(request dto.SendNotificationPaymentRequest) (string, error) {
	var notificationMsg string
	var err error

	user, err := s.userService.GetUsersById(request.UserID)
	if err != nil {
		return "", err
	}
	orders, err := s.repo.GetOrderById(request.OrderID)
	if err != nil {
		return "", err
	}

	switch request.PaymentStatus {
	case "Menunggu Konfirmasi":
		notificationMsg = fmt.Sprintf("Alloo, %s! Pesananmu dengan ID %s udah berhasil dibuat, nih. Ditunggu yupp!!", user.Name, orders.IdOrder)
	case "Konfirmasi":
		notificationMsg = fmt.Sprintf("Thengkyuu, %s! Pembayaran untuk pesananmu dengan ID %s udah kami terima, nih. Semoga harimu menyenangkan!", user.Name, orders.IdOrder)
	case "Gagal":
		notificationMsg = fmt.Sprintf("Maaf, %s. Pembayaran untuk pesanan dengan ID %s gagal, nih. Beritahu kami apabila kamu butuh bantuan yaa!!", user.Name, orders.IdOrder)
	default:
		return "", errors.New("Status pesanan tidak valid")
	}

	notificationRequest := sendnotif.SendNotificationRequest{
		OrderID: request.OrderID,
		UserID:  request.UserID,
		Title:   "Status Pembayaran",
		Body:    notificationMsg,
		Token:   user.DeviceToken,
	}
	_, _, err = s.fcmService.CreateFcm(notificationRequest)
	if err != nil {
		logrus.Error("Gagal mengirim notifikasi: ", err)
		return "", err
	}

	return notificationMsg, nil
}

func (s *OrderService) SendNotificationOrder(request dto.SendNotificationOrderRequest) (string, error) {
	var notificationMsg string
	var err error

	user, err := s.userService.GetUsersById(request.UserID)
	if err != nil {
		return "", err
	}
	orders, err := s.repo.GetOrderById(request.OrderID)
	if err != nil {
		return "", err
	}

	switch request.OrderStatus {
	case "Pengiriman":
		notificationMsg = fmt.Sprintf("Alloo, %s! Pesanan dengan ID %s udah dalam proses pengiriman, nih. Mohon ditunggu yupp!", user.Name, orders.IdOrder)
	case "Selesai":
		notificationMsg = fmt.Sprintf("Yeayy, %s! Pesananmu dengan ID %s udah sampai tujuan, nih. Semoga sukakk yupp!", user.Name, orders.IdOrder)
	case "Menunggu Konfirmasi":
		notificationMsg = fmt.Sprintf("Alloo, %s! Pesananmu dengan ID %s sedang menunggu konfirmasi, nih. Ditunggu yupp!", user.Name, orders.IdOrder)
	case "Proses":
		notificationMsg = fmt.Sprintf("Alloo, %s! Pesananmu dengan ID %s sedang dalam proses, nih. Ditunggu yupp!", user.Name, orders.IdOrder)
	case "Gagal":
		notificationMsg = fmt.Sprintf("Sowwy, %s. Pesananmu dengan ID %s gagal. Coba lagi, yukk!", user.Name, orders.IdOrder)
	default:
		return "", errors.New("Status pengiriman tidak valid")
	}

	notificationRequest := sendnotif.SendNotificationRequest{
		OrderID: request.OrderID,
		UserID:  request.UserID,
		Title:   "Status Pengiriman",
		Body:    notificationMsg,
		Token:   user.DeviceToken,
	}
	_, _, err = s.fcmService.CreateFcm(notificationRequest)
	if err != nil {
		logrus.Error("Gagal mengirim notifikasi: ", err)
		return "", err
	}

	return notificationMsg, nil
}

func (s *OrderService) GetFilterDateRange(filterType string) (time.Time, time.Time, error) {
	filterType = strings.ToLower(filterType)
	now := time.Now()

	switch filterType {
	case "minggu ini":
		startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
		startDate := time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 0, 7)
		return startDate, endDate, nil
	case "bulan ini":
		startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		nextMonth := startDate.AddDate(0, 1, 0)
		endDate := nextMonth.Add(-time.Second)
		return startDate, endDate, nil
	case "tahun ini":
		startDate := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		nextYear := startDate.AddDate(1, 0, 0)
		endDate := nextYear.Add(-time.Second)
		return startDate, endDate, nil
	default:
		return time.Time{}, time.Time{}, errors.New("tipe filter tidak valid")
	}
}
