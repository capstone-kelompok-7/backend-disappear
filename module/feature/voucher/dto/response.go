package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	_ "text/template/parse"
	"time"
)

type VoucherModelsResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Category  string `json:"category"`
	Discouunt int    `json:"discount"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"stop_date"`
	Status    string `json:"status"`
}

type VoucherModelsResponseAll struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Category  string `json:"category"`
	Discouunt int    `json:"discount"`
	StartDate string `json:"start_date"`
	EndDatee  string `json:"stop_date"`
	Status    string `json:"status"`
}

func VoucherResponseFormatterCreate(voucher *entities.VoucherModels) VoucherModelsResponse {
	voucherFormatter := VoucherModelsResponse{}
	voucherFormatter.ID = voucher.ID
	voucherFormatter.Name = voucher.Name
	voucherFormatter.Code = voucher.Code
	voucherFormatter.Category = voucher.Category
	voucherFormatter.Discouunt = voucher.Discouunt

	sampleFormat := "2006-01-02"

	var dateNow = time.Now()
	dateNowFormat := dateNow.Format(sampleFormat)

	noww, _ := time.Parse("2006-01-02", dateNowFormat)
	endd, _ := time.Parse("2006-01-02", voucher.EndDate)

	if noww.After(endd) {
		voucherFormatter.Status = "kadaluarsa"
	} else if noww.Before(endd) {
		voucherFormatter.Status = "aktif"
	}

	voucherFormatter.StartDate = voucher.StartDate
	voucherFormatter.EndDate = voucher.EndDate

	return voucherFormatter
}
func VoucherResponseFormatter(voucher *entities.VoucherModels) VoucherModelsResponse {
	voucherFormatter := VoucherModelsResponse{}
	voucherFormatter.ID = voucher.ID
	voucherFormatter.Name = voucher.Name
	voucherFormatter.Code = voucher.Code
	voucherFormatter.Category = voucher.Category
	voucherFormatter.Discouunt = voucher.Discouunt

	sampleFormat := "2006-01-02"

	var dateNow = time.Now()
	dateNowFormat := dateNow.Format(sampleFormat)

	parsedStartDate, _ := time.Parse(time.RFC3339Nano, voucher.StartDate)
	formatstartdate := parsedStartDate.Format("2006-01-02")
	voucherFormatter.StartDate = formatstartdate

	parsedenddate, _ := time.Parse(time.RFC3339Nano, voucher.EndDate)
	formatenddate := parsedenddate.Format(sampleFormat)
	voucherFormatter.EndDate = formatenddate

	noww, _ := time.Parse("2006-01-02", dateNowFormat)
	endd, _ := time.Parse("2006-01-02", formatenddate)

	if noww.After(endd) {
		voucherFormatter.Status = "kadaluarsa"
	} else if noww.Before(endd) {
		voucherFormatter.Status = "aktif"
	}

	voucherFormatter.StartDate = formatstartdate
	voucherFormatter.EndDate = formatenddate

	return voucherFormatter
}
func VoucherResponseFormatterAll(voucher entities.VoucherModels) VoucherModelsResponseAll {
	voucherFormatter := VoucherModelsResponseAll{}
	voucherFormatter.ID = voucher.ID
	voucherFormatter.Name = voucher.Name
	voucherFormatter.Code = voucher.Code
	voucherFormatter.Category = voucher.Category
	voucherFormatter.Discouunt = voucher.Discouunt
	voucherFormatter.StartDate = voucher.StartDate

	return voucherFormatter
}

func VoucherModelsFormatterAll(vouchers []entities.VoucherModels) []VoucherModelsResponseAll {
	var voucherFormatter []VoucherModelsResponseAll

	for _, voucher := range vouchers {
		formatVoucher := VoucherResponseFormatterAll(voucher)

		var dateNow = time.Now()
		formatteddatenow := dateNow.Format("2006-01-02")

		parsedStartDate, _ := time.Parse(time.RFC3339Nano, voucher.StartDate)
		formatstartdate := parsedStartDate.Format("2006-01-02")
		formatVoucher.StartDate = formatstartdate

		parsedEndDate, _ := time.Parse(time.RFC3339Nano, voucher.EndDate)
		formatenddate := parsedEndDate.Format("2006-01-02")
		formatVoucher.EndDatee = formatenddate

		noww, _ := time.Parse("2006-01-02", formatteddatenow)
		endd, _ := time.Parse("2006-01-02", formatenddate)

		if noww.After(endd) {
			formatVoucher.Status = "kadaluarsa"
		} else if noww.Before(endd) {
			formatVoucher.Status = "aktif"
		}

		voucherFormatter = append(voucherFormatter, formatVoucher)

	}

	return voucherFormatter
}
