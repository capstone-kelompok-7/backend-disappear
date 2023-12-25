package repository

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils/binderbyte"
	"github.com/capstone-kelompok-7/backend-disappear/utils/payment"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type OrderRepository struct {
	db         *gorm.DB
	coreClient coreapi.Client
}

func NewOrderRepository(db *gorm.DB, coreClient coreapi.Client) order.RepositoryOrderInterface {
	return &OrderRepository{
		db:         db,
		coreClient: coreClient,
	}
}

func (r *OrderRepository) FindByName(page, perPage int, name string) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	offset := (page - 1) * perPage
	query := r.db.Preload("User").Offset(offset).Limit(perPage).
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.deleted_at IS NULL")

	if name != "" {
		query = query.Where("users.name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetTotalCustomerCountByName(name string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.OrderModels{}).
		Joins("JOIN users ON orders.user_id = users.id").
		Where("users.deleted_at IS NULL")

	if name != "" {
		query = query.Where("users.name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *OrderRepository) FindAll(page, perPage int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Preload("User").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetTotalOrderCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.OrderModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *OrderRepository) GetOrderById(orderID string) (*entities.OrderModels, error) {
	var orders entities.OrderModels

	if err := r.db.
		Preload("OrderDetails").
		Preload("OrderDetails.Product").
		Preload("OrderDetails.Product.ProductPhotos").
		Preload("User").
		Preload("Voucher").
		Preload("Address").
		Where("id = ? AND deleted_at IS NULL", orderID).
		First(&orders).
		Error; err != nil {
		return nil, err
	}

	return &orders, nil
}

func (r *OrderRepository) CreateOrder(newOrder *entities.OrderModels) (*entities.OrderModels, error) {
	err := r.db.Create(newOrder).Error
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (r *OrderRepository) CreateOrderDetails(newOrderDetails *entities.OrderDetailsModels) (*entities.OrderDetailsModels, error) {
	err := r.db.Create(newOrderDetails).Error
	if err != nil {
		return nil, err
	}
	return newOrderDetails, nil
}

func (r *OrderRepository) ConfirmPayment(orderID, orderStatus, paymentStatus string) error {
	var orders entities.OrderModels
	if err := r.db.Model(&orders).Where("id = ?", orderID).Updates(map[string]interface{}{
		"order_status":   orderStatus,
		"payment_status": paymentStatus,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) ProcessGatewayPayment(totalAmountPaid uint64, orderID string, paymentMethod, name, email string) (interface{}, error) {
	var paymentType coreapi.CoreapiPaymentType

	switch paymentMethod {
	case "qris":
		paymentType = coreapi.PaymentTypeQris
	case "bank_transfer":
		paymentType = coreapi.PaymentTypeBankTransfer
	case "gopay":
		paymentType = coreapi.PaymentTypeGopay
	}

	coreClient := r.coreClient
	resp, err := payment.CreateCoreAPIPaymentRequest(coreClient, orderID, int64(totalAmountPaid), paymentType, name, email)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("gagal membuat permintaan pembayaran")
	}

	return resp, nil
}

func (r *OrderRepository) CheckTransaction(orderID string) (dto.Status, error) {
	var status dto.Status
	transactionStatusResp, err := r.coreClient.CheckTransaction(orderID)
	if err != nil {
		return dto.Status{}, err
	} else {
		if transactionStatusResp != nil {
			status = payment.TransactionStatus(transactionStatusResp)
			return status, nil
		}
	}
	return dto.Status{}, err
}

func (r *OrderRepository) UpdateOrderStatus(req *dto.UpdateOrderStatus) error {
	orders := entities.OrderModels{}
	if err := r.db.Where("id = ?", req.OrderID).First(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("pesanan tidak ditemukan")
		}
		return err
	}

	if err := r.db.Model(&orders).Updates(map[string]interface{}{
		"status_order_date": req.StatusOrderDate,
		"order_status":      req.OrderStatus,
		"extra_info":        req.ExtraInfo,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetAllOrdersByUserID(userID uint64) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels

	if err := r.db.
		Preload("OrderDetails").
		Preload("OrderDetails.Product").
		Preload("OrderDetails.Product.ProductPhotos").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetAllOrdersWithFilter(userID uint64, orderStatus string) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels

	if err := r.db.
		Preload("OrderDetails").
		Preload("OrderDetails.Product").
		Preload("OrderDetails.Product.ProductPhotos").
		Where("user_id = ? AND order_status = ? AND deleted_at IS NULL", userID, orderStatus).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) AcceptOrder(orderID, orderStatus string) error {
	if err := r.db.Model(&entities.OrderModels{}).
		Where("id = ?", orderID).
		Update("order_status", orderStatus).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) Tracking(courier, awb string) (map[string]interface{}, error) {
	result, err := binderbyte.TrackingPackages(courier, awb)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepository) GetOrderByDateRange(startDate, endDate time.Time, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.Preload("User").Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountByDateRange(startDate, endDate time.Time) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrderByOrderStatus(orderStatus string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.Preload("User").Where("order_status = ? AND deleted_at IS NULL", orderStatus).Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountByByOrderStatus(orderStatus string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Where("order_status = ? AND deleted_at IS NULL", orderStatus).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrderByDateRangeAndStatus(startDate, endDate time.Time, status string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.Preload("User").Where("created_at BETWEEN ? AND ? AND order_status = ? AND deleted_at IS NULL", startDate, endDate, status).
		Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountByDateRangeAndStatus(startDate, endDate time.Time, status string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Where("created_at BETWEEN ? AND ? AND order_status = ? AND deleted_at IS NULL", startDate, endDate, status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrderByDateRangeAndStatusAndSearch(startDate, endDate time.Time, status, search string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.Preload("User").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.order_status = ? AND users.name LIKE ? AND orders.deleted_at IS NULL",
			startDate, endDate, status, "%"+search+"%").
		Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountByDateRangeAndStatusAndSearch(startDate, endDate time.Time, status, search string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.order_status = ? AND users.name LIKE ? AND orders.deleted_at IS NULL",
			startDate, endDate, status, "%"+search+"%").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrdersBySearchAndDateRange(startDate, endDate time.Time, search string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.
		Preload("User").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.created_at BETWEEN ? AND ? AND users.name LIKE ? AND orders.deleted_at IS NULL", startDate, endDate, "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountBySearchAndDateRange(startDate, endDate time.Time, search string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.created_at BETWEEN ? AND ? AND users.name LIKE ? AND orders.deleted_at IS NULL", startDate, endDate, "%"+search+"%").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrdersBySearchAndStatus(status, search string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.
		Preload("User").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("users.name LIKE ? AND orders.order_status = ? AND orders.deleted_at IS NULL", "%"+search+"%", status).
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrdersCountBySearchAndStatus(status, search string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Joins("JOIN users ON users.id = orders.user_id").
		Where("users.name LIKE ? AND orders.order_status = ? AND orders.deleted_at IS NULL", "%"+search+"%", status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrderBySearchAndPaymentStatus(status, search string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.
		Preload("User").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("users.name LIKE ? AND orders.payment_status = ? AND orders.deleted_at IS NULL", "%"+search+"%", status).
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountBySearchAndPaymentStatus(status, search string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Joins("JOIN users ON users.id = orders.user_id").
		Where("users.name LIKE ? AND orders.payment_status = ? AND orders.deleted_at IS NULL", "%"+search+"%", status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrderByDateRangeAndPaymentStatusAndSearch(startDate, endDate time.Time, status, search string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.Preload("User").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.payment_status = ? AND users.name LIKE ? AND orders.deleted_at IS NULL",
			startDate, endDate, status, "%"+search+"%").
		Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountByDateRangeAndPaymentStatusAndSearch(startDate, endDate time.Time, status, search string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.payment_status = ? AND users.name LIKE ? AND orders.deleted_at IS NULL",
			startDate, endDate, status, "%"+search+"%").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrderByPaymentStatus(orderStatus string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.Preload("User").Where("payment_status = ? AND deleted_at IS NULL", orderStatus).Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountByByPaymentStatus(orderStatus string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Where("payment_status = ? AND deleted_at IS NULL", orderStatus).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetOrderByDateRangeAndPaymentStatus(startDate, endDate time.Time, status string, offset, limit int) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels
	if err := r.db.Preload("User").Where("created_at BETWEEN ? AND ? AND payment_status = ? AND deleted_at IS NULL", startDate, endDate, status).
		Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderCountByDateRangeAndPaymentStatus(startDate, endDate time.Time, status string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Where("created_at BETWEEN ? AND ? AND payment_status = ? AND deleted_at IS NULL", startDate, endDate, status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
