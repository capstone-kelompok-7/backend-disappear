package fcm

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/utils/sendnotif"
	"github.com/labstack/echo/v4"
)

type RepositoryFcmInterface interface {
	CreateFcm(fcm *entities.FcmModels) (*entities.FcmModels, error)
	GetFcmByIdUser(id uint64) ([]*entities.FcmModels, error)
	GetFcmById(id uint64) (*entities.FcmModels, error)
	DeleteFcmById(id uint64) error
	SendMessageNotification(request sendnotif.SendNotificationRequest) (string, error)
}
type ServiceFcmInterface interface {
	CreateFcm(request sendnotif.SendNotificationRequest) (string, *entities.FcmModels, error)
	GetFcmByIdUser(id uint64) ([]*entities.FcmModels, error)
	GetFcmById(id uint64) (*entities.FcmModels, error)
	DeleteFcmById(id uint64) error
}

type HandlerFcmInterface interface {
	CreateFcm() echo.HandlerFunc
	GetFcmByIdUser() echo.HandlerFunc
	GetFcmById() echo.HandlerFunc
	DeleteFcmById() echo.HandlerFunc
}
