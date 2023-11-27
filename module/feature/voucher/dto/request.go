package dto

import "time"

type CreateVoucherRequest struct {
	Name        string    `form:"name" json:"name" validate:"required"`
	Code        string    `form:"code" json:"code" validate:"required"`
	Category    string    `form:"category" json:"category" validate:"required"`
	Description string    `form:"description" json:"description" validate:"required"`
	Discount    uint64    `form:"discount" json:"discount" validate:"required"`
	StartDate   time.Time `form:"start_date" json:"start_date" validate:"required"`
	EndDate     time.Time `form:"end_date" json:"end_date" validate:"required"`
	MinPurchase uint64    `form:"min_purchase" json:"min_purchase" validate:"required"`
	Stock       uint64    `form:"stock" json:"stock" validate:"required"`
	Status      string    `form:"status" json:"status"`
}

type UpdateVoucherRequest struct {
	Name        string    `form:"name" json:"name"`
	Code        string    `form:"code" json:"code"`
	Category    string    `form:"category" json:"category"`
	Description string    `form:"description" json:"description"`
	Discount    uint64    `form:"discount" json:"discount"`
	StartDate   time.Time `form:"start_date" json:"start_date"`
	EndDate     time.Time `form:"end_date" json:"end_date"`
	MinPurchase uint64    `form:"min_purchase" json:"min_purchase"`
	Stock       uint64    `form:"stock" json:"stock"`
	Status      string    `form:"status" json:"status"`
}

type ClaimsVoucherRequest struct {
	UserID    uint64 `form:"user_id" json:"user_id"`
	VoucherID uint64 `form:"voucher_id" json:"voucher_id" validate:"required"`
}
