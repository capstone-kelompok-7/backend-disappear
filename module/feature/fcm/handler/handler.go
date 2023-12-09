package handler

import (
	"strconv"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type FcmHandler struct {
	service fcm.ServiceFcmInterface
}

func NewFcmHandler(service fcm.ServiceFcmInterface) fcm.HandlerFcmInterface {
	return &FcmHandler{
		service: service,
	}
}

func (h *FcmHandler) CreateFcm() echo.HandlerFunc {
	return func(c echo.Context) error {
		fcmRequest := new(dto.FcmRequest)
		if err := c.Bind(fcmRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai: "+err.Error())
		}

		if err := utils.ValidateStruct(fcmRequest); err != nil {
			return response.SendBadRequestResponse(c, "Validasi gagal: "+err.Error())
		}

		newFcm := &entities.FcmModels{
			Title: fcmRequest.Title,
			Body:  fcmRequest.Body,
		}

		res, statusFcm, err := h.service.CreateFcm(newFcm, fcmRequest.Token)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mengirim noifikasi: "+err.Error())
		}

		return response.SendStatusCreatedResponse(c, "Berhasil mengirim noifikasi", dto.FormatFcmCreate(res, statusFcm, fcmRequest.Token))

	}
}
func (h *FcmHandler) GetFcmByIdUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var fcm []*entities.FcmModels
		var err error

		idUser := c.Param("id")
		idUserConf, _ := strconv.Atoi(idUser)

		fcm, err = h.service.GetFcmByIdUser(uint64(idUserConf))
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar notifikasi by iduser: "+err.Error())
		}

		return response.SendStatusOkWithDataResponse(c, "Berhasil mendapatkan daftar notifikasi by id user", dto.FormatFcmGetbyIdUser2(fcm))
	}
}
func (h *FcmHandler) GetFcmById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var fcm *entities.FcmModels
		var err error

		id := c.Param("id")
		idConf, _ := strconv.Atoi(id)

		fcm, err = h.service.GetFcmById(uint64(idConf))
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar notifikasi by id: "+err.Error())
		}

		return response.SendStatusOkWithDataResponse(c, "Berhasil mendapatkan notifikasi by id", dto.FormatFcmGetbyIdUser(fcm))

	}
}
func (h *FcmHandler) DeleteFcmById() echo.HandlerFunc {
	return func(c echo.Context) error {

		var err error

		id := c.Param("id")
		idConf, _ := strconv.Atoi(id)

		err = h.service.DeleteFcmById(uint64(idConf))
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan daftar notifikasi by id: "+err.Error())
		}

		return response.SendStatusOkWithDataResponse(c, "Berhasil menghapus notifikasi by id", nil)
	}
}
