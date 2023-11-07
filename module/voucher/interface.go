package voucher

import (
	voucher "github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
	"github.com/labstack/echo/v4"
)

type RepositoryVoucherInterface interface {
	CreateVoucher(newData voucher.VoucherModels) (*voucher.VoucherModels, error)
	GetAllVouchers(page int, limit int, search string) ([]voucher.VoucherModels, error)
	DeleteVoucherById(id int) error
	// DeleteVoucherByName(name string) error
	EditVoucherById(data voucher.VoucherModels) (*voucher.VoucherModels, error)
}

type ServiceVoucherInterface interface {
	CreateVoucher(newData voucher.VoucherModels) (*voucher.VoucherModels, error)
	GetAllVouchers(page int, limit int, search string) ([]voucher.VoucherModels, error)
	DeleteVoucherById(id int) error
	// DeleteVoucherByName(name string) error
	EditVoucherById(data voucher.VoucherModels) (*voucher.VoucherModels, error)
}

type HandlerVoucherInterface interface {
	CreateVoucher() echo.HandlerFunc
	GetAllVouchers() echo.HandlerFunc
	DeleteVoucherById() echo.HandlerFunc
	// DeleteVoucherByName() echo.HandlerFunc
	EditVoucherById() echo.HandlerFunc
}
