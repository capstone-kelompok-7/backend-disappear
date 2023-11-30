package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type AddressResponse struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	AcceptedName string `json:"accepted_name" `
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	IsPrimary    bool   `json:"is_primary"`
}

func FormatAddress(address *entities.AddressModels) *AddressResponse {
	addressFormatter := &AddressResponse{
		ID:           address.ID,
		UserID:       address.UserID,
		AcceptedName: address.AcceptedName,
		Phone:        address.Phone,
		Address:      address.Address,
		IsPrimary:    address.IsPrimary,
	}
	return addressFormatter
}

func FormatterAddress(addresses []*entities.AddressModels) []*AddressResponse {
	addressFormatters := make([]*AddressResponse, 0)

	for _, address := range addresses {
		formattedAddress := FormatAddress(address)
		addressFormatters = append(addressFormatters, formattedAddress)
	}

	return addressFormatters
}
