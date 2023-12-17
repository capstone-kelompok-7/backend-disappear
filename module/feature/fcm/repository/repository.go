package repository

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/utils/sendnotif"
	"github.com/sirupsen/logrus"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm"

	"gorm.io/gorm"
)

type FcmRepository struct {
	db  *gorm.DB
	fcm sendnotif.FcmServiceInterface
}

func NewFcmRepository(db *gorm.DB, fcm sendnotif.FcmServiceInterface) fcm.RepositoryFcmInterface {
	return &FcmRepository{
		db:  db,
		fcm: fcm,
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

func (r *FcmRepository) SendMessageNotification(request sendnotif.SendNotificationRequest) (string, error) {
	var err error
	var sendSuccess string

	sendSuccess, err = r.fcm.SendNotification(request)
	if err != nil {
		logrus.Error("Failed to send notification")
		return "", err
	}

	return sendSuccess, nil
}
