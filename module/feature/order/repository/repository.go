package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/order"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) order.RepositoryOrderInterface {
	return &OrderRepository{
		db: db,
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
