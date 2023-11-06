package voucher

import (
	voucher "github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryVoucherInterface interface {
	CreateVoucher(newData voucher.VoucherModels) (*voucher.VoucherModels, error)
	GetAllVouchers(page int, limit int) ([]voucher.VoucherModels, error)
	// GetVoucherByName(name string) (*voucher.VoucherModels, error)
	// DeleteVoucherByName(name string) error
	// EditVoucherByName(name string) (*voucher.VoucherModels, error)
}

type ServiceVoucherInterface interface {
	CreateVoucher(newData voucher.VoucherModels) (*voucher.VoucherModels, error)
	GetAllVouchers(page int, limit int) ([]voucher.VoucherModels, error)
	// GetVoucherByName(name string) (*voucher.VoucherModels, error)
	// DeleteVoucherByName(name string) error
	// EditVoucherByName(name string) (*voucher.VoucherModels, error)
}

type HandlerVoucherInterface interface {
	CreateVoucher() echo.HandlerFunc
	GetAllVouchers() echo.HandlerFunc
	// GetVoucherByName() echo.HandlerFunc
	// DeleteVoucherByName() echo.HandlerFunc
	// EditVoucherByName() echo.HandlerFunc
}
