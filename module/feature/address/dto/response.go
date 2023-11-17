package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type AddressFormatter struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	AcceptedName string `json:"accepted_name"`
	Street       string `json:"street"`
	SubDistrict  string `json:"sub_district"`
	City         string `json:"city"`
	Province     string `json:"province"`
	PostalCode   int    `json:"postal_code"`
	Note         string `json:"note"`
}

func FormatAddress(address *entities.AddressModels) *AddressFormatter {
	addressFormatter := &AddressFormatter{}
	addressFormatter.ID = address.ID
	addressFormatter.UserID = address.UserID
	addressFormatter.AcceptedName = address.AcceptedName
	addressFormatter.Street = address.Street
	addressFormatter.SubDistrict = address.SubDistrict
	addressFormatter.City = address.City
	addressFormatter.Province = address.Province
	addressFormatter.PostalCode = address.PostalCode
	addressFormatter.Note = address.Note

	return addressFormatter
}

func FormatterAddress(addresses []*entities.AddressModels) []*AddressFormatter {
	addressFormatters := make([]*AddressFormatter, 0)

	for _, address := range addresses {
		formattedAddress := FormatAddress(address)
		addressFormatters = append(addressFormatters, formattedAddress)
	}

	return addressFormatters
}
