package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"time"
)

type CreateOrderRequest struct {
	AddressID     uint64 `form:"address_id" json:"address_id" validate:"required"`
	VoucherID     uint64 `form:"voucher_id" json:"voucher_id"`
	Note          string `form:"note" json:"note"`
	ProductID     uint64 `json:"product_id" validate:"required"`
	Quantity      uint64 `json:"quantity" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
}

type CreateOrderCartRequest struct {
	AddressID     uint64                    `form:"address_id" json:"address_id" validate:"required"`
	VoucherID     uint64                    `form:"voucher_id" json:"voucher_id"`
	Note          string                    `form:"note" json:"note"`
	PaymentMethod string                    `json:"payment_method" validate:"required"`
	CartItems     []entities.CartItemModels `json:"cart_items" validate:"required"`
}

type PaymentOrderRequest struct {
	OrderID     string `json:"order_id"`
	PaymentType string `json:"payment_type"`
	TotalAmount int64  `json:"total_amount"`
}

type Status struct {
	PaymentStatus string
	OrderStatus   string
}

type UpdateOrderStatus struct {
	OrderID         string    `json:"order_id"`
	StatusOrderDate time.Time `json:"status_order_date"`
	OrderStatus     string    `json:"order_status"`
	ExtraInfo       string    `json:"extra_info"`
}

type SendNotificationPaymentRequest struct {
	PaymentStatus string `json:"payment_status"`
	OrderID       string `json:"order_id"`
	UserID        uint64 `json:"user_id"`
	UserName      string `json:"user_name"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	Token         string `json:"token"`
}

type SendNotificationOrderRequest struct {
	OrderStatus string `json:"order_status"`
	OrderID     string `json:"order_id"`
	UserID      uint64 `json:"user_id"`
	UserName    string `json:"user_name"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Token       string `json:"token"`
}
