package service

import (
	"errors"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm"
	"github.com/capstone-kelompok-7/backend-disappear/utils/sendnotif"
	"github.com/sirupsen/logrus"
)

type FcmService struct {
	repo fcm.RepositoryFcmInterface
}

func NewFcmService(repo fcm.RepositoryFcmInterface) fcm.ServiceFcmInterface {
	return &FcmService{
		repo: repo,
	}
}

func (s *FcmService) CreateFcm(request sendnotif.SendNotificationRequest) (string, *entities.FcmModels, error) {
	var sendSuccess string
	var err error

	sendSuccess, err = s.repo.SendMessageNotification(request)
	if err != nil {
		logrus.Error("Failed to send notification:", err)
		return "", nil, err
	}

	value := &entities.FcmModels{
		OrderID: request.OrderID,
		UserID:  request.UserID,
		Title:   request.Title,
		Body:    request.Body,
	}

	response, createErr := s.repo.CreateFcm(value)
	if createErr != nil {
		logrus.Error("Failed to create notification in the database:", createErr)
		return "", nil, createErr
	}

	return sendSuccess, response, nil
}

func (s *FcmService) GetFcmByIdUser(id uint64) ([]*entities.FcmModels, error) {
	res, err := s.repo.GetFcmByIdUser(id)
	if err != nil {
		return nil, errors.New("Notifikasi tidak ditemukan")
	}
	return res, nil
}

func (s *FcmService) GetFcmById(id uint64) (*entities.FcmModels, error) {
	res, err := s.repo.GetFcmById(id)
	if err != nil {
		return nil, errors.New("Notifikasi tidak ditemukan")
	}
	return res, nil
}

func (s *FcmService) DeleteFcmById(id uint64) error {
	result, err := s.repo.GetFcmById(id)
	if err != nil {
		return errors.New("fcm tidak ditemukan")
	}

	err = s.repo.DeleteFcmById(result.ID)
	if err != nil {
		return errors.New("fcm id tidak ditemukan")
	}

	return nil
}
