package voucher

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/dto"
	"github.com/labstack/echo/v4"
)

type RepositoryVoucherInterface interface {
	CreateVoucher(newData entities.VoucherModels) (entities.VoucherModels, error)
	FindVoucherByName(page, perPage int, name string) ([]entities.VoucherModels, error)
	GetTotalVoucherCountByName(name string) (int64, error)
	FindAllVoucher(page, perPage int) ([]entities.VoucherModels, error)
	GetTotalVoucherCount() (int64, error)
	DeleteVoucher(id uint64) error
	GetVoucherById(id uint64) (entities.VoucherModels, error)
	UpdateVoucher(id uint64, updatedVoucher dto.UpdateVoucherRequest) (entities.VoucherModels, error)
}

type ServiceVoucherInterface interface {
	CreateVoucher(newData entities.VoucherModels) (entities.VoucherModels, error)
	DeleteVoucher(id uint64) error
	GetVoucherById(id uint64) (entities.VoucherModels, error)
	UpdateVoucher(id uint64, updatedData dto.UpdateVoucherRequest) (entities.VoucherModels, error)
	GetAllVoucher(page, perPage int) ([]entities.VoucherModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetVouchersByName(page, perPage int, name string) ([]entities.VoucherModels, int64, error)
}

type HandlerVoucherInterface interface {
	CreateVoucher() echo.HandlerFunc
	GetAllVouchers() echo.HandlerFunc
	DeleteVoucherById() echo.HandlerFunc
	GetVoucherById() echo.HandlerFunc
	UpdateVouchers() echo.HandlerFunc
}
