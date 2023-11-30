package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type CreateOrderRequest struct {
	AddressID uint64 `form:"address_id" json:"address_id" validate:"required"`
	VoucherID uint64 `form:"voucher_id" json:"voucher_id"`
	Note      string `form:"note" json:"note"`
	ProductID uint64 `json:"product_id" validate:"required"`
	Quantity  uint64 `json:"quantity" validate:"required"`
}

type CreateOrderCartRequest struct {
	AddressID uint64                    `form:"address_id" json:"address_id" validate:"required"`
	VoucherID uint64                    `form:"voucher_id" json:"voucher_id"`
	Note      string                    `form:"note" json:"note"`
	CartItems []entities.CartItemModels `json:"cart_items" validate:"required"`
}
