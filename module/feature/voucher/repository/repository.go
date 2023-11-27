package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
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

func (r *VoucherRepository) CreateVoucher(newData *entities.VoucherModels) (*entities.VoucherModels, error) {
	if err := r.db.Create(&newData).Error; err != nil {
		return newData, err
	}

	return newData, nil
}

func (r *VoucherRepository) UpdateVoucher(voucherID uint64, updatedVoucher *entities.VoucherModels) error {
	var vouchers *entities.CarouselModels
	if err := r.db.Where("id = ? AND deleted_at IS NULL", voucherID).First(&vouchers).Error; err != nil {
		return err
	}
	if err := r.db.Updates(&updatedVoucher).Error; err != nil {
		return err
	}
	return nil
}

func (r *VoucherRepository) DeleteVoucher(voucherID uint64) error {
	vouchers := &entities.VoucherModels{}
	if err := r.db.First(vouchers, voucherID).Error; err != nil {
		return err
	}

	if err := r.db.Model(vouchers).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil

}

func (r *VoucherRepository) GetVoucherById(voucherID uint64) (*entities.VoucherModels, error) {
	var vouchers *entities.VoucherModels
	if err := r.db.Where("id = ? AND deleted_at IS NULL", voucherID).First(&vouchers).Error; err != nil {
		return nil, err
	}
	return vouchers, nil
}

func (r *VoucherRepository) FindAllVoucher(page, perPage int) ([]*entities.VoucherModels, error) {
	var vouchers []*entities.VoucherModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Where("deleted_at IS NULL").Find(&vouchers).Error
	if err != nil {
		return nil, err
	}
	return vouchers, nil
}

func (r *VoucherRepository) GetTotalVoucherCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.VoucherModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *VoucherRepository) IsVoucherAlreadyClaimed(userID uint64, voucherID uint64) (bool, error) {
	var existingClaim entities.VoucherClaimModels
	err := r.db.Where("user_id = ? AND voucher_id = ?", userID, voucherID).First(&existingClaim).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *VoucherRepository) ClaimVoucher(claimVoucher *entities.VoucherClaimModels) error {
	if err := r.db.Create(&claimVoucher).Error; err != nil {
		return err
	}
	return nil
}

func (r *VoucherRepository) ReduceStockWhenClaimed(voucherID, quantity uint64) error {
	var claims entities.VoucherModels
	if err := r.db.Model(&claims).Where("id = ?", voucherID).Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
		return err
	}
	return nil
}

func (r *VoucherRepository) GetVoucherCategory(voucherID uint64) (string, error) {
	var vouchers entities.VoucherModels
	if err := r.db.Select("category").Where("id = ? AND deleted_at IS NULL", voucherID).First(&vouchers).Error; err != nil {
		return "", err
	}
	return vouchers.Category, nil
}

func (r *VoucherRepository) DeleteUserVoucherClaims(userID, voucherID uint64) error {
	if err := r.db.Where("user_id = ? AND voucher_id = ?", userID, voucherID).Delete(&entities.VoucherClaimModels{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *VoucherRepository) GetUserVoucherClaims(userID uint64) ([]*entities.VoucherClaimModels, error) {
	var userVouchers []*entities.VoucherClaimModels

	if err := r.db.Preload("Voucher").Where("user_id = ?", userID).Find(&userVouchers).Error; err != nil {
		return nil, err
	}

	return userVouchers, nil
}

func (r *VoucherRepository) GetVoucherByCode(code string) (*entities.VoucherModels, error) {
	var vouchers entities.VoucherModels
	if err := r.db.Where("code = ? AND deleted_at IS NULL", code).First(&vouchers).Error; err != nil {
		return nil, err
	}
	return &vouchers, nil
}
