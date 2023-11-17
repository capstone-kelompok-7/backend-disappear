package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) address.RepositoryAddressInterface {
	return &AddressRepository{
		db: db,
	}
}

func (r *AddressRepository) CreateAddress(newData *entities.AddressModels) (*entities.AddressModels, error) {
	if err := r.db.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (r *AddressRepository) FindAllByUserID(userID uint64, page, perPage int) ([]*entities.AddressModels, error) {
	var addresses []*entities.AddressModels
	offset := (page - 1) * perPage
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Limit(perPage).
		Offset(offset).
		Find(&addresses).
		Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepository) GetTotalAddressCountByUserID(userID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&entities.AddressModels{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).
		Error
	return count, err
}
