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
	query := r.db.Offset(offset).Limit(perPage).
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

func (r *OrderRepository) ProcessGatewayPayment(totalAmountPaid uint64, orderID string, paymentMethod string) (interface{}, error) {
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
	resp, err := payment.CreateCoreAPIPaymentRequest(coreClient, orderID, int64(totalAmountPaid), paymentType)
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
