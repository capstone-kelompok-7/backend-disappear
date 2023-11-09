package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
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

func (r *VoucherRepository) CreateVoucher(newData entities.VoucherModels) (*entities.VoucherModels, error) {
	if err := r.db.Create(&newData).Error; err != nil {
		return nil, err
	}

	return &newData, nil
}
func (r *VoucherRepository) GetAllVouchers(currentPage int, limit int, search string) ([]entities.VoucherModels, error) {
	var listVoucher = []entities.VoucherModels{}

	if search != "" {
		if err := r.db.Where("category LIKE ?", "%"+search+"%").Offset((currentPage - 1) * limit).Limit(limit).Find(&listVoucher).Error; err != nil {
			return nil, err
		}
	} else if search == "" {
		if err := r.db.Offset((currentPage - 1) * limit).Limit(limit).Find(&listVoucher).Error; err != nil {
			return nil, err
		}
	}

	return listVoucher, nil
}

func (r *VoucherRepository) GetAllVouchersToCalculatePage() ([]entities.VoucherModels, error) {
	var listVoucher = []entities.VoucherModels{}

	if err := r.db.Find(&listVoucher).Error; err != nil {
		return nil, err
	}

	return listVoucher, nil
}

func (r *VoucherRepository) EditVoucherById(data entities.VoucherModels) (*entities.VoucherModels, error) {
	var voucher = entities.VoucherModels{}
	if err := r.db.Model(&voucher).Where("id = ?", data.ID).Updates(map[string]interface{}{
		"name":        data.Name,
		"code":        data.Code,
		"category":    data.Category,
		"description": data.Description,
		"discount":    data.Discouunt,
		"start_date":  data.StartDate,
		"end_date":    data.EndDate,
		"min_amount":  data.MinAmount,
	}).Error; err != nil {
		return nil, err
	}

	return &voucher, nil
}

func (r *VoucherRepository) DeleteVoucherById(id int) error {
	var voucher = entities.VoucherModels{}

	if err := r.db.Unscoped().Where("id = ?", id).Delete(&voucher).Error; err != nil {
		return nil
	}
	return nil

}

func (r *VoucherRepository) GetVoucherById(id int) (*entities.VoucherModels, error) {
	var voucher = entities.VoucherModels{}
	if err := r.db.Where("id = ?", id).First(&voucher).Error; err != nil {
		return nil, err
	}

	return &voucher, nil
}
