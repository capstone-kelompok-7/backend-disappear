package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address/dto"
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
		Phone:        addressData.Phone,
		Address:      addressData.Address,
		IsPrimary:    addressData.IsPrimary,
		CreatedAt:    time.Now(),
	}

	createdAddress, err := s.repo.CreateAddress(value)
	if err != nil {
		return nil, err
	}
	if createdAddress.IsPrimary {
		currentPrimaryAddress, err := s.repo.GetPrimaryAddressByUserID(createdAddress.UserID)
		if err != nil {
			return nil, errors.New("gagal mendapatkan alamat utama")
		}
		if currentPrimaryAddress != nil && currentPrimaryAddress.ID != createdAddress.ID {
			err = s.repo.UpdateIsPrimary(currentPrimaryAddress.ID, false)
			if err != nil {
				return nil, errors.New("gagal merubah alamat utama")
			}
		}
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

func (s *AddressService) GetAddressByID(addressID uint64) (*entities.AddressModels, error) {
	result, err := s.repo.GetAddressByID(addressID)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan")
	}
	return result, nil
}

func (s *AddressService) UpdateAddress(userID, addressID uint64, updatedAddress *dto.UpdateAddressRequest) error {
	existingAddress, err := s.repo.GetAddressByID(addressID)
	if err != nil {
		return errors.New("alamat tidak ditemukan")
	}

	if updatedAddress.IsPrimary {
		err := s.repo.UpdateIsPrimary(addressID, true)
		if err != nil {
			return errors.New("gagal merubah alamat utama")
		}

		currentPrimaryAddress, err := s.repo.GetPrimaryAddressByUserID(userID)
		if err == nil && currentPrimaryAddress != nil && currentPrimaryAddress.ID != existingAddress.ID {
			err = s.repo.UpdateIsPrimary(currentPrimaryAddress.ID, false)
			if err != nil {
				return errors.New("gagal menonaktifkan alamat utama sebelumnya")
			}
		}
	}

	err = s.repo.UpdateAddress(addressID, updatedAddress)
	if err != nil {
		return errors.New("gagal memperbarui alamat")
	}

	return nil
}

func (s *AddressService) DeleteAddress(addressID, userID uint64) error {
	addresses, err := s.repo.GetAddressByID(addressID)
	if err != nil {
		return errors.New("Alamat tidak ditemukan")
	}
	primaryAddress, err := s.repo.GetPrimaryAddressByUserID(userID)
	if err != nil {
		return errors.New("Gagal mendapatkan alamat utama")
	}

	if addresses.ID == primaryAddress.ID {
		return errors.New("Alamat utama tidak dapat dihapus")
	}

	err = s.repo.DeleteAddress(addresses.ID)
	if err != nil {
		return errors.New("Gagal menghapus alamat")
	}

	return nil
}
