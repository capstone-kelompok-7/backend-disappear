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

func (s *FcmService) CreateFcm(fcm *entities.FcmModels, token string) (*entities.FcmModels, string, error) {
	value := &entities.FcmModels{
		UserID: fcm.UserID,
		Title:  fcm.Title,
		Body:   fcm.Body,
	}
	var err error
	var response *entities.FcmModels
	var sendsuccess string

	sendsuccess, err = sendnotif.SendNotification(fcm.Title, fcm.Body, token)
	if err == nil {
		response, err = s.repo.CreateFcm(value)
		if err != nil {
			logrus.Error("gagal membuat notification ke database")
		}
	} else if err != nil {
		logrus.Error("gagal mengirim notification")
	}

	return response, sendsuccess, nil

}

func (s *FcmService) GetFcmByIdUser(id uint64) ([]*entities.FcmModels, error) {
	res, err := s.repo.GetFcmByIdUser(id)
	if err != nil {
		return nil, errors.New("Notifikasi tidak ditemukan" + err.Error())
	}
	return res, nil
}

func (s *FcmService) GetFcmById(id uint64) (*entities.FcmModels, error) {
	res, err := s.repo.GetFcmById(id)
	if err != nil {
		return nil, errors.New("Notifikasi tidak ditemukan" + err.Error())
	}
	return res, nil
}

func (s *FcmService) DeleteFcmById(id uint64) error {
	findfcm, err := s.repo.GetFcmById(id)
	if err != nil {
		return errors.New("fcm tidak ditemukan: " + err.Error())
	}

	if findfcm == nil {
		return errors.New("fcm tidak ditemukan: " + err.Error())
	}

	err = s.repo.DeleteFcmById(id)
	if err != nil {
		return errors.New("gagal menghapus artikel: " + err.Error())
	}

	return nil
}
