package dto

import (
	_ "text/template/parse"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
)

type VoucherFormatter struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Discount    uint64    `json:"discount"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end-date"`
	MinPurchase uint64    `json:"min_purchase" `
	Stock       uint64    `json:"stock" `
	Status      string    `json:"status" `
}

func FormatVoucher(voucher *entities.VoucherModels) *VoucherFormatter {
	voucherFormatter := &VoucherFormatter{}
	voucherFormatter.ID = voucher.ID
	voucherFormatter.Name = voucher.Name
	voucherFormatter.Code = voucher.Code
	voucherFormatter.Category = voucher.Category
	voucherFormatter.Description = voucher.Description
	voucherFormatter.Discount = voucher.Discount
	voucherFormatter.StartDate = voucher.StartDate
	voucherFormatter.EndDate = voucher.EndDate
	voucherFormatter.MinPurchase = voucher.MinPurchase
	voucherFormatter.Stock = voucher.Stock
	voucherFormatter.Status = voucher.Status

	return voucherFormatter
}

func FormatterVoucher(vouchers []*entities.VoucherModels) []*VoucherFormatter {
	var voucherFormatters []*VoucherFormatter

	for _, voucher := range vouchers {
		formattedVoucher := FormatVoucher(voucher)
		voucherFormatters = append(voucherFormatters, formattedVoucher)
	}

	return voucherFormatters
}

type GetVoucherUserVoucherFormatter struct {
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Discount    uint64    `json:"discount"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	MinPurchase uint64    `json:"min_purchase" `
	Status      string    `json:"status" `
}

type GetVoucherUserResponse struct {
	ID        uint64                         `json:"id"`
	UserID    uint64                         `json:"user_id"`
	VoucherID uint64                         `json:"voucher_id"`
	Voucher   GetVoucherUserVoucherFormatter `json:"voucher"`
}

func GetVoucherUserFormatter(voucherClaims []*entities.VoucherClaimModels) ([]GetVoucherUserResponse, error) {
	var voucherResponses []GetVoucherUserResponse

	for _, voucher := range voucherClaims {
		voucherResponse := GetVoucherUserResponse{
			ID:        voucher.ID,
			UserID:    voucher.UserID,
			VoucherID: voucher.VoucherID,
			Voucher: GetVoucherUserVoucherFormatter{
				Name:        voucher.Voucher.Name,
				Code:        voucher.Voucher.Code,
				Category:    voucher.Voucher.Category,
				Description: voucher.Voucher.Description,
				Discount:    voucher.Voucher.Discount,
				StartDate:   voucher.Voucher.StartDate,
				EndDate:     voucher.Voucher.EndDate,
				MinPurchase: voucher.Voucher.MinPurchase,
				Status:      voucher.Voucher.Status,
			},
		}
		voucherResponses = append(voucherResponses, voucherResponse)
	}

	return voucherResponses, nil
}

type VoucherToClaimsFormatter struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Discount    uint64    `json:"discount"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end-date"`
	MinPurchase uint64    `json:"min_purchase" `
	Stock       uint64    `json:"stock" `
}

func FormatVoucherToClaims(voucher *entities.VoucherModels) *VoucherToClaimsFormatter {
	voucherToClaimsFormatter := &VoucherToClaimsFormatter{}
	voucherToClaimsFormatter.ID = voucher.ID
	voucherToClaimsFormatter.Name = voucher.Name
	voucherToClaimsFormatter.Code = voucher.Code
	voucherToClaimsFormatter.Category = voucher.Category
	voucherToClaimsFormatter.Description = voucher.Description
	voucherToClaimsFormatter.Discount = voucher.Discount
	voucherToClaimsFormatter.StartDate = voucher.StartDate
	voucherToClaimsFormatter.EndDate = voucher.EndDate
	voucherToClaimsFormatter.MinPurchase = voucher.MinPurchase
	voucherToClaimsFormatter.Stock = voucher.Stock

	return voucherToClaimsFormatter
}

func FormatterVoucherToClaims(vouchers []*entities.VoucherModels) []*VoucherToClaimsFormatter {
	var voucherToClaimsFormatters []*VoucherToClaimsFormatter

	for _, voucher := range vouchers {
		formattedVoucher := FormatVoucherToClaims(voucher)
		voucherToClaimsFormatters = append(voucherToClaimsFormatters, formattedVoucher)
	}

	return voucherToClaimsFormatters
}
