package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type FcmFormatterCreate struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Token  string `json:"token"`
	Status string `json:"status"`
}

type FcmFormatterByIdUser struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func FormatFcmCreate(fcm *entities.FcmModels, status string, token string) *FcmFormatterCreate {
	FcmFormatter := &FcmFormatterCreate{}
	FcmFormatter.ID = fcm.ID
	FcmFormatter.UserID = fcm.UserID
	FcmFormatter.Title = fcm.Title
	FcmFormatter.Body = fcm.Body
	FcmFormatter.Status = status
	FcmFormatter.Token = token

	return FcmFormatter
}
func FormatFcmGetbyIdUser(fcm *entities.FcmModels) *FcmFormatterByIdUser {
	FcmFormatter := &FcmFormatterByIdUser{}
	FcmFormatter.ID = fcm.ID
	FcmFormatter.UserID = fcm.UserID
	FcmFormatter.Title = fcm.Title
	FcmFormatter.Body = fcm.Body

	return FcmFormatter
}

func FormatFcmGetbyIdUser2(fcms []*entities.FcmModels) []*FcmFormatterByIdUser {
	var FcmFormatter []*FcmFormatterByIdUser

	for _, fcm := range fcms {
		formatFcm := FormatFcmGetbyIdUser(fcm)
		FcmFormatter = append(FcmFormatter, formatFcm)
	}

	return FcmFormatter
}
