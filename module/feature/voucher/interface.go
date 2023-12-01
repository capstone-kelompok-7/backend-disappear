package voucher

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryVoucherInterface interface {
	CreateVoucher(newData *entities.VoucherModels) (*entities.VoucherModels, error)
	FindAllVoucher(page, perPage int) ([]*entities.VoucherModels, error)
	GetTotalVoucherCount() (int64, error)
	DeleteVoucher(voucherID uint64) error
	GetVoucherById(voucherID uint64) (*entities.VoucherModels, error)
	UpdateVoucher(voucherID uint64, updatedVoucher *entities.VoucherModels) error
	IsVoucherAlreadyClaimed(userID uint64, voucherID uint64) (bool, error)
	ClaimVoucher(claimVoucher *entities.VoucherClaimModels) error
	ReduceStockWhenClaimed(voucherID, quantity uint64) error
	GetVoucherCategory(voucherID uint64) (string, error)
	DeleteUserVoucherClaims(userID, voucherID uint64) error
	GetUserVoucherClaims(userID uint64) ([]*entities.VoucherClaimModels, error)
	GetVoucherByCode(code string) (*entities.VoucherModels, error)
	FindByStatus(page, perPage int, status string) ([]*entities.VoucherModels, error)
	GetTotalVoucherCountByStatus(status string) (int64, error)
	FindByCategory(page, perPage int, category string) ([]*entities.VoucherModels, error)
	GetTotalVoucherCountByCategory(category string) (int64, error)
	FindByStatusCategory(page, perPage int, status, category string) ([]*entities.VoucherModels, error)
	GetTotalVoucherCountByStatusCategory(status, category string) (int64, error)
	FindAllVoucherToClaims(limit int, userID uint64) ([]*entities.VoucherModels, error)
}

type ServiceVoucherInterface interface {
	CreateVoucher(newData *entities.VoucherModels) (*entities.VoucherModels, error)
	DeleteVoucher(voucherID uint64) error
	GetVoucherById(voucherID uint64) (*entities.VoucherModels, error)
	UpdateVoucher(voucherID uint64, req *entities.VoucherModels) error
	GetAllVoucher(page, perPage int) ([]*entities.VoucherModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	CanClaimsVoucher(userID, voucherID uint64) (bool, error)
	ClaimVoucher(req *entities.VoucherClaimModels) error
	DeleteVoucherClaims(userID, voucherID uint64) error
	GetUserVouchers(userID uint64) ([]*entities.VoucherClaimModels, error)
	GetVoucherByStatus(page, perPage int, status string) ([]*entities.VoucherModels, int64, error)
	GetVoucherByCategory(page, perPage int, category string) ([]*entities.VoucherModels, int64, error)
	GetVoucherByStatusCategory(page, perPage int, status, category string) ([]*entities.VoucherModels, int64, error)
	GetAllVoucherToClaims(limit int, userID uint64) ([]*entities.VoucherModels, error)
}

type HandlerVoucherInterface interface {
	CreateVoucher() echo.HandlerFunc
	GetAllVouchers() echo.HandlerFunc
	DeleteVoucherById() echo.HandlerFunc
	GetVoucherById() echo.HandlerFunc
	UpdateVouchers() echo.HandlerFunc
	ClaimVoucher() echo.HandlerFunc
	GetVoucherUser() echo.HandlerFunc
	GetAllVouchersToClaims() echo.HandlerFunc
}
