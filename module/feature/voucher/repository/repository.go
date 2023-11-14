package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/dto"
	"gorm.io/gorm"
	"time"
)

type VoucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) voucher.RepositoryVoucherInterface {
	return &VoucherRepository{
		db: db,
	}
}

func (r *VoucherRepository) CreateVoucher(newData entities.VoucherModels) (entities.VoucherModels, error) {
	if err := r.db.Create(&newData).Error; err != nil {
		return newData, err
	}

	return newData, nil
}

func (r *VoucherRepository) UpdateVoucher(id uint64, updatedVoucher dto.UpdateVoucherRequest) (entities.VoucherModels, error) {
	var vouchers entities.VoucherModels
	if err := r.db.Model(&entities.VoucherModels{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updatedVoucher).Error; err != nil {
		return vouchers, err
	}
	return vouchers, nil
}

func (r *VoucherRepository) DeleteVoucher(id uint64) error {
	vouchers := &entities.VoucherModels{}
	if err := r.db.First(vouchers, id).Error; err != nil {
		return err
	}

	if err := r.db.Model(vouchers).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil

}

func (r *VoucherRepository) GetVoucherById(id uint64) (entities.VoucherModels, error) {
	var voucher = entities.VoucherModels{}
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&voucher).Error; err != nil {
		return voucher, err
	}

	return voucher, nil
}

func (r *VoucherRepository) FindVoucherByName(page, perPage int, name string) ([]entities.VoucherModels, error) {
	var vouchers []entities.VoucherModels
	offset := (page - 1) * perPage
	query := r.db.Offset(offset).Limit(perPage).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&vouchers).Error
	if err != nil {
		return vouchers, err
	}

	return vouchers, nil
}

func (r *VoucherRepository) GetTotalVoucherCountByName(name string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.VoucherModels{}).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *VoucherRepository) FindAllVoucher(page, perPage int) ([]entities.VoucherModels, error) {
	var vouchers []entities.VoucherModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Where("deleted_at IS NULL").Find(&vouchers).Error
	if err != nil {
		return vouchers, err
	}
	return vouchers, nil
}

func (r *VoucherRepository) GetTotalVoucherCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.VoucherModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}
