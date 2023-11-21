package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address/dto"
	"gorm.io/gorm"
	"time"
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

func (r *AddressRepository) GetAddressByID(addressID uint64) (*entities.AddressModels, error) {
	var addresses entities.AddressModels
	err := r.db.Where("id = ? AND deleted_at IS NULL", addressID).First(&addresses).Error
	if err != nil {
		return nil, err
	}
	return &addresses, nil
}

func (r *AddressRepository) UpdateAddress(addressID uint64, updatedAddress *dto.UpdateAddressRequest) error {
	err := r.db.Model(&entities.AddressModels{}).
		Where("id = ?", addressID).
		Updates(updatedAddress).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AddressRepository) GetPrimaryAddressByUserID(userID uint64) (*entities.AddressModels, error) {
	var addresses *entities.AddressModels
	err := r.db.Where("user_id = ? AND is_primary = ? AND deleted_at IS NULL", userID, true).First(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepository) UpdateIsPrimary(addressID uint64, isPrimary bool) error {
	var addresses *entities.AddressModels
	err := r.db.Model(&addresses).Where("id = ?", addressID).Update("is_primary", isPrimary).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AddressRepository) DeleteAddress(addressID uint64) error {
	addresses := &entities.AddressModels{}
	if err := r.db.First(&addresses, addressID).Error; err != nil {
		return err
	}

	if err := r.db.Model(&addresses).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}
