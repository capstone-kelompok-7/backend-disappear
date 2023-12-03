package service

import (
	"errors"
	"testing"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateService_CreateAddress(t *testing.T) {
	repo := mocks.NewRepositoryAddressInterface(t)
	service := NewAddressService(repo)

	currentPrimaryAddress := &entities.AddressModels{
		ID:           1,
		UserID:       1,
		AcceptedName: "test",
		Phone:        "0812345678910",
		Address:      "Jl. Test",
		IsPrimary:    true,
		CreatedAt:    time.Now(),
	}

	expectedAddress := &entities.AddressModels{
		ID:           2,
		UserID:       currentPrimaryAddress.UserID,
		AcceptedName: currentPrimaryAddress.AcceptedName,
		Phone:        currentPrimaryAddress.Phone,
		Address:      "Jl. Test Dulu",
		IsPrimary:    true,
		CreatedAt:    time.Now().AddDate(0, 0, 7),
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("CreateAddress", mock.AnythingOfType("*entities.AddressModels")).Return(expectedAddress, nil).Once()
		repo.On("GetPrimaryAddressByUserID", mock.AnythingOfType("uint64")).Return(currentPrimaryAddress, nil).Once()
		repo.On("UpdateIsPrimary", currentPrimaryAddress.ID, false).Return(nil).Once()

		result, err := service.CreateAddress(expectedAddress)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedAddress, result)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CreateAddress", func(t *testing.T) {
		expectedError := errors.New("gagal mendapatkan alamat utama")
		repo.On("CreateAddress", mock.AnythingOfType("*entities.AddressModels")).Return(nil, expectedError).Once()

		result, err := service.CreateAddress(expectedAddress)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetPrimaryAddressByUserID", func(t *testing.T) {
		expectedError := errors.New("gagal mendapatkan alamat utama")
		repo.On("CreateAddress", mock.AnythingOfType("*entities.AddressModels")).Return(expectedAddress, nil).Once()
		repo.On("GetPrimaryAddressByUserID", mock.AnythingOfType("uint64")).Return(nil, expectedError).Once()

		result, err := service.CreateAddress(expectedAddress)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - UpdateIsPrimary", func(t *testing.T) {
		expectedError := errors.New("gagal merubah alamat utama")
		repo.On("CreateAddress", mock.AnythingOfType("*entities.AddressModels")).Return(expectedAddress, nil).Once()
		repo.On("GetPrimaryAddressByUserID", mock.AnythingOfType("uint64")).Return(currentPrimaryAddress, nil).Once()
		repo.On("UpdateIsPrimary", currentPrimaryAddress.ID, false).Return(expectedError).Once()

		result, err := service.CreateAddress(expectedAddress)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)

		repo.AssertExpectations(t)
	})
}

func TestAddressService_GetAll(t *testing.T) {
	repo := mocks.NewRepositoryAddressInterface(t)
	service := NewAddressService(repo)

	address := &entities.AddressModels{
		ID:           1,
		UserID:       1,
		AcceptedName: "test",
		Phone:        "123456789010",
		Address:      "Jl. Test",
		IsPrimary:    true,
		CreatedAt:    time.Now(),
	}

	t.Run("Success Case - Address Found", func(t *testing.T) {
		repo.On("FindAllByUserID", uint64(1), 1, 10).Return([]*entities.AddressModels{address}, nil).Once()
		repo.On("GetTotalAddressCountByUserID", uint64(1)).Return(int64(1), nil).Once()

		result, totalItems, err := service.GetAll(uint64(1), 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, []*entities.AddressModels{address}, result)
		assert.Equal(t, int64(1), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - FindAllByUserID fails", func(t *testing.T) {
		expectedError := errors.New("failed to find addresses")
		repo.On("FindAllByUserID", uint64(1), 1, 10).Return(nil, expectedError).Once()

		result, totalItems, err := service.GetAll(uint64(1), 1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedError, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetTotalAddressCountByUserID Error", func(t *testing.T) {
		expectedErr := errors.New("GetTotalAddressCountByUserID Error")
		repo.On("FindAllByUserID", uint64(1), 1, 10).Return([]*entities.AddressModels{address}, nil).Once()
		repo.On("GetTotalAddressCountByUserID", uint64(1)).Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetAll(uint64(1), 1, 10)

		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestAddressService_CalculatePaginationValues(t *testing.T) {
	service := &AddressService{}

	t.Run("Page less than or equal to zero should default to 1", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(0, 100, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page exceeds total pages should set to total pages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(15, 100, 8)

		assert.Equal(t, 13, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page within limits should return correct values", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(2, 100, 8)

		assert.Equal(t, 2, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Total items not perfectly divisible by perPage should round totalPages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(1, 95, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 12, totalPages)
	})
}

func TestAddressService_GetNextPage(t *testing.T) {
	service := &AddressService{}

	t.Run("Next Page Within Total Pages", func(t *testing.T) {
		currentPage := 3
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, currentPage+1, nextPage)
	})

	t.Run("Next Page Equal to Total Pages", func(t *testing.T) {
		currentPage := 5
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, totalPages, nextPage)
	})
}

func TestAddressService_GetPrevPage(t *testing.T) {
	service := &AddressService{}

	t.Run("Previous Page Within Bounds", func(t *testing.T) {
		currentPage := 3

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage-1, prevPage)
	})

	t.Run("Previous Page at Lower Bound", func(t *testing.T) {
		currentPage := 1

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage, prevPage)
	})
}

func TestAddressService_GetAddresById(t *testing.T) {
	repo := mocks.NewRepositoryAddressInterface(t)
	service := NewAddressService(repo)

	address := &entities.AddressModels{
		ID:           1,
		UserID:       1,
		AcceptedName: "test",
		Phone:        "123456789010",
		Address:      "Jl. Test",
		IsPrimary:    true,
		CreatedAt:    time.Now(),
	}

	t.Run("Success Case - Address Found", func(t *testing.T) {
		addressID := uint64(1)
		repo.On("GetAddressByID", addressID).Return(address, nil).Once()

		result, err := service.GetAddressByID(addressID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, address.ID, result.ID)
		assert.Equal(t, address.UserID, result.UserID)
		assert.Equal(t, address.AcceptedName, result.AcceptedName)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Adress Not Found", func(t *testing.T) {
		addressID := uint64(2)

		expectedErr := errors.New("alamat tidak ditemukan")
		repo.On("GetAddressByID", addressID).Return(nil, expectedErr).Once()

		result, err := service.GetAddressByID(addressID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})
}

func TestAddressService_UpdateAddress(t *testing.T) {
	repo := mocks.NewRepositoryAddressInterface(t)
	service := NewAddressService(repo)

	userID := uint64(1)
	addressID := uint64(1)

	existingAddress := &entities.AddressModels{
		ID:           addressID,
		UserID:       userID,
		AcceptedName: "test",
		Phone:        "08122345678",
		Address:      "Jl. Test",
		IsPrimary:    false,
		CreatedAt:    time.Now(),
	}

	updatedAddress := &dto.UpdateAddressRequest{
		AcceptedName: "updated",
		Phone:        "08123456789",
		Address:      "Jl. Updated",
		IsPrimary:    true,
	}

	t.Run("Success Case - Update Primary Address", func(t *testing.T) {
		repo.On("GetAddressByID", addressID).Return(existingAddress, nil).Once()
		repo.On("UpdateIsPrimary", addressID, true).Return(nil).Once()
		repo.On("GetPrimaryAddressByUserID", userID).Return(nil, nil).Once()
		repo.On("UpdateAddress", addressID, updatedAddress).Return(nil).Once()

		result := service.UpdateAddress(userID, addressID, updatedAddress)

		assert.NoError(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - UpdateIsPrimary Error", func(t *testing.T) {
		repo.On("GetAddressByID", addressID).Return(existingAddress, nil).Once()
		repo.On("UpdateIsPrimary", addressID, true).Return(errors.New("gagal merubah alamat utama")).Once()

		err := service.UpdateAddress(userID, addressID, updatedAddress)

		assert.Error(t, err)
		assert.Equal(t, "gagal merubah alamat utama", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - GetAddressByID Error", func(t *testing.T) {
		expectedErr := errors.New("GetAddressByID error")
		repo.On("GetAddressByID", addressID).Return(nil, expectedErr).Once()

		err := service.UpdateAddress(userID, addressID, updatedAddress)

		assert.Error(t, err)
		assert.Equal(t, "alamat tidak ditemukan", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - Disable Previous Primary Address Error", func(t *testing.T) {
		repo.On("GetAddressByID", addressID).Return(existingAddress, nil).Once()
		repo.On("UpdateIsPrimary", addressID, true).Return(nil).Once()

		previousPrimaryAddress := &entities.AddressModels{
			ID:           2,
			UserID:       userID,
			AcceptedName: "previous_primary",
			Phone:        "08123456789",
			Address:      "Jl. Previous Primary",
			IsPrimary:    true,
			CreatedAt:    time.Now(),
		}
		repo.On("GetPrimaryAddressByUserID", userID).Return(previousPrimaryAddress, nil).Once()

		repo.On("UpdateIsPrimary", previousPrimaryAddress.ID, false).Return(errors.New("gagal menonaktifkan alamat utama sebelumnya")).Once()

		err := service.UpdateAddress(userID, addressID, updatedAddress)

		assert.Error(t, err)
		assert.Equal(t, "gagal menonaktifkan alamat utama sebelumnya", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failure Case - UpdateAddress Error", func(t *testing.T) {
		repo.On("GetAddressByID", addressID).Return(existingAddress, nil).Once()
		repo.On("UpdateIsPrimary", addressID, true).Return(nil).Once()
		repo.On("GetPrimaryAddressByUserID", userID).Return(nil, nil).Once()

		repo.On("UpdateAddress", addressID, updatedAddress).Return(errors.New("gagal memperbarui alamat")).Once()

		err := service.UpdateAddress(userID, addressID, updatedAddress)

		assert.Error(t, err)
		assert.Equal(t, "gagal memperbarui alamat", err.Error())
		repo.AssertExpectations(t)
	})
}

func TestAddressService_DeleteAddress(t *testing.T) {
	repo := mocks.NewRepositoryAddressInterface(t)
	service := NewAddressService(repo)

	userID := uint64(1)
	addressID := uint64(2)

	address := &entities.AddressModels{
		ID:           addressID,
		UserID:       userID,
		AcceptedName: "test",
		Phone:        "08122345678",
		Address:      "Jl. Test",
		IsPrimary:    false,
		CreatedAt:    time.Now(),
	}

	primaryAddress := &entities.AddressModels{
		ID:           1,
		UserID:       userID,
		AcceptedName: "primary",
		Phone:        "08123456789",
		Address:      "Jl. Primary",
		IsPrimary:    true,
		CreatedAt:    time.Now(),
	}

	t.Run("Success Case - Delete Address", func(t *testing.T) {
		repo.On("GetAddressByID", addressID).Return(address, nil).Once()
		repo.On("GetPrimaryAddressByUserID", userID).Return(primaryAddress, nil).Once()
		repo.On("DeleteAddress", addressID).Return(nil).Once()

		err := service.DeleteAddress(addressID, userID)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Address Not Found", func(t *testing.T) {
		repo.On("GetAddressByID", addressID).Return(nil, errors.New("Alamat tidak ditemukan")).Once()

		err := service.DeleteAddress(addressID, userID)

		assert.Error(t, err)
		assert.Equal(t, "Alamat tidak ditemukan", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Get Primary Address Error", func(t *testing.T) {
		repo.On("GetAddressByID", addressID).Return(address, nil).Once()
		repo.On("GetPrimaryAddressByUserID", userID).Return(nil, errors.New("Gagal mendapatkan alamat utama")).Once()

		err := service.DeleteAddress(addressID, userID)

		assert.Error(t, err)
		assert.Equal(t, "Gagal mendapatkan alamat utama", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Attempt to Delete Primary Address", func(t *testing.T) {
		repo.On("GetAddressByID", primaryAddress.ID).Return(primaryAddress, nil).Once()
		repo.On("GetPrimaryAddressByUserID", userID).Return(primaryAddress, nil).Once()

		err := service.DeleteAddress(primaryAddress.ID, userID)

		assert.Error(t, err)
		assert.Equal(t, "Alamat utama tidak dapat dihapus", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Delete Address Error", func(t *testing.T) {
		repo.On("GetAddressByID", addressID).Return(address, nil).Once()
		repo.On("GetPrimaryAddressByUserID", userID).Return(primaryAddress, nil).Once()
		repo.On("DeleteAddress", addressID).Return(errors.New("Gagal menghapus alamat")).Once()

		err := service.DeleteAddress(addressID, userID)

		assert.Error(t, err)
		assert.Equal(t, "Gagal menghapus alamat", err.Error())
		repo.AssertExpectations(t)
	})
}
