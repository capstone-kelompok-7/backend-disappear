package repository

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm"

	"gorm.io/gorm"
)

type FcmRepository struct {
	db *gorm.DB
}

func NewFcmRepository(db *gorm.DB) fcm.RepositoryFcmInterface {
	return &FcmRepository{
		db: db,
	}
}

func (r *FcmRepository) CreateFcm(fcm *entities.FcmModels) (*entities.FcmModels, error) {
	if err := r.db.Create(fcm).Error; err != nil {
		return nil, err
	}

	return fcm, nil
}

func (r *FcmRepository) GetFcmByIdUser(id uint64) ([]*entities.FcmModels, error) {
	var fcm []*entities.FcmModels
	err := r.db.Where("deleted_at IS NULL && user_id = ?", id).Find(&fcm).Error
	if err != nil {
		return nil, err
	}

	return fcm, nil
}

func (r *FcmRepository) GetFcmById(id uint64) (*entities.FcmModels, error) {
	var fcm entities.FcmModels

	if err := r.db.Where("deleted_at IS NULL && id = ?", id).First(&fcm).Error; err != nil {
		return nil, err
	}

	return &fcm, nil
}

func (r *FcmRepository) DeleteFcmById(id uint64) error {
	var fcm = &entities.FcmModels{}

	if err := r.db.First(fcm, id).Error; err != nil {
		return err
	}

	if err := r.db.Model(fcm).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil

}
