package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	_ "text/template/parse"
	"time"
)

type VoucherFormatter struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Category    string    `json:"category"`
	Discount    uint64    `json:"discount"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"stop_date"`
	MinPurchase uint64    `json:"min_purchase" `
	Stock       uint64    `json:"stock" `
	Status      string    `json:"status" `
}

func FormatVoucher(voucher entities.VoucherModels) VoucherFormatter {
	voucherFormatter := VoucherFormatter{}
	voucherFormatter.ID = voucher.ID
	voucherFormatter.Name = voucher.Name
	voucherFormatter.Code = voucher.Code
	voucherFormatter.Category = voucher.Category
	voucherFormatter.Discount = voucher.Discount
	voucherFormatter.StartDate = voucher.StartDate
	voucherFormatter.EndDate = voucher.EndDate
	voucherFormatter.MinPurchase = voucher.MinPurchase
	voucherFormatter.Stock = voucher.Stock
	voucherFormatter.Status = voucher.Status

	return voucherFormatter
}

func FormatterVoucher(vouchers []entities.VoucherModels) []VoucherFormatter {
	var voucherFormatters []VoucherFormatter

	for _, voucher := range vouchers {
		formattedVoucher := FormatVoucher(voucher)
		voucherFormatters = append(voucherFormatters, formattedVoucher)
	}

	return voucherFormatters
}
