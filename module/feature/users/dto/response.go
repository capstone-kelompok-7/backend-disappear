package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
)

type UserFormatter struct {
	ID           uint64             `json:"id"`
	Email        string             `json:"email"`
	Role         string             `json:"role"`
	Name         string             `json:"name"`
	Phone        string             `json:"phone"`
	PhotoProfile string             `json:"photo_profile"`
	TotalGram    uint64             `json:"total_gram"`
	IsVerified   bool               `json:"is_verified"`
	Level        string             `json:"level"`
	Exp          uint64             `json:"exp"`
	Addresses    []AddressFormatter `json:"addresses"`
}

type AddressFormatter struct {
	ID           uint64 `json:"id"`
	AcceptedName string `json:"accepted_name"`
	Street       string `json:"street"`
	SubDistrict  string `json:"sub_district"`
	City         string `json:"city"`
	Province     string `json:"province"`
	PostalCode   int    `json:"postal_code"`
	Note         string `json:"note"`
}

func FormatUser(user *entities.UserModels) *UserFormatter {
	userFormatter := &UserFormatter{}
	userFormatter.ID = user.ID
	userFormatter.Email = user.Email
	userFormatter.Role = user.Role
	userFormatter.Name = user.Name
	userFormatter.Phone = user.Phone
	userFormatter.PhotoProfile = user.PhotoProfile
	userFormatter.TotalGram = user.TotalGram
	userFormatter.IsVerified = user.IsVerified
	userFormatter.Level = user.Level
	userFormatter.Exp = user.Exp

	var addresses []AddressFormatter
	for _, address := range user.Address {
		addressesFormatter := AddressFormatter{
			ID:           address.ID,
			AcceptedName: address.AcceptedName,
			Street:       address.Street,
			SubDistrict:  address.SubDistrict,
			City:         address.City,
			Province:     address.Province,
			PostalCode:   address.PostalCode,
			Note:         address.Note,
		}
		addresses = append(addresses, addressesFormatter)
	}
	userFormatter.Addresses = addresses

	return userFormatter
}

func FormatterUsers(users []*entities.UserModels) []*UserFormatter {
	usersFormatters := make([]*UserFormatter, 0)

	for _, user := range users {
		formattedUsers := FormatUser(user)
		usersFormatters = append(usersFormatters, formattedUsers)
	}

	return usersFormatters
}
