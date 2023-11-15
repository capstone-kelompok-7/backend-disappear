package dto

import "time"

type CreateVoucherRequest struct {
	Name        string    `json:"name" validate:"required"`
	Code        string    `json:"code" validate:"required"`
	Category    string    `json:"category" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Discount    uint64    `json:"discount" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	MinPurchase uint64    `json:"min_purchase" validate:"required"`
	Stock       uint64    `json:"stock" validate:"required"`
	Status      string    `json:"status"`
}

type UpdateVoucherRequest struct {
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Category    string    `json:"category" `
	Description string    `json:"description"`
	Discount    uint64    `json:"discount"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	MinPurchase uint64    `json:"min_purchase" `
	Stock       uint64    `json:"stock" `
	Status      string    `json:"status"`
}
