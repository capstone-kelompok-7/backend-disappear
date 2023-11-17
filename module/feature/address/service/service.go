package service

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"math"
	"time"
)

type AddressService struct {
	repo address.RepositoryAddressInterface
}

func NewAddressService(addressRepo address.RepositoryAddressInterface) address.ServiceAddressInterface {
	return &AddressService{
		repo: addressRepo,
	}
}

func (s *AddressService) CreateAddress(addressData *entities.AddressModels) (*entities.AddressModels, error) {
	value := &entities.AddressModels{
		UserID:       addressData.UserID,
		AcceptedName: addressData.AcceptedName,
		Street:       addressData.Street,
		SubDistrict:  addressData.SubDistrict,
		City:         addressData.City,
		Province:     addressData.Province,
		PostalCode:   addressData.PostalCode,
		Note:         addressData.Note,
		CreatedAt:    time.Now(),
	}

	createdAddress, err := s.repo.CreateAddress(value)
	if err != nil {
		return nil, err
	}

	return createdAddress, nil
}

func (s *AddressService) GetAll(userID uint64, page, perPage int) ([]*entities.AddressModels, int64, error) {
	addresses, err := s.repo.FindAllByUserID(userID, page, perPage)
	if err != nil {
		return addresses, 0, err
	}

	totalItems, err := s.repo.GetTotalAddressCountByUserID(userID)
	if err != nil {
		return addresses, 0, err
	}

	return addresses, totalItems, nil
}

func (s *AddressService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	pageInt := page
	if pageInt <= 0 {
		pageInt = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	if pageInt > totalPages {
		pageInt = totalPages
	}

	return pageInt, totalPages
}

func (s *AddressService) GetNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}
	return totalPages
}

func (s *AddressService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}
	return 1
}
