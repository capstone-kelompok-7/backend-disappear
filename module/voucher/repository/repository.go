package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
	"gorm.io/gorm"
)

type VoucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) voucher.RepositoryVoucherInterface {
	return &VoucherRepository{
		db: db,
	}
}

func (r *VoucherRepository) CreateVoucher(newData domain.VoucherModels) (*domain.VoucherModels, error) {
	if err := r.db.Create(&newData).Error; err != nil {
		return nil, err
	}

	return &newData, nil
}
func (r *VoucherRepository) GetAllVouchers(currentPage int, limit int) ([]domain.VoucherModels, error) {
	var listVoucher = []domain.VoucherModels{}

	if err := r.db.Offset((currentPage - 1) * limit).Limit(limit).Find(&listVoucher).Error; err != nil {
		return nil, err
	}

	return listVoucher, nil
}

// func (r *VoucherRepository) GetAllVouchers(page int, limit int) ([]domain.VoucherModels, error) {

// }

// func (r *VoucherRepository) GetVoucherByName(name string) (*domain.VoucherModels, error) {

// }

// func (r *VoucherRepository) EditVoucherByName(name string) (*domain.VoucherModels, error) {

// }

// func (r *VoucherRepository) DeleteVoucherByName(name string) (*domain.VoucherModels, error) {

// }
