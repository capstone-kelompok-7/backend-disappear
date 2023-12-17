package service

import (
	"errors"
	"testing"
	"time"
	_ "time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	_ "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	_ "github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/dto"
	user_mock "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/mocks"
	user_service "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/service"
	voucherMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/mocks"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	_ "github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "github.com/stretchr/testify/mock"
)

func TestVocucherService_calculatepaginations(t *testing.T) {
	service := &VoucherService{}

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

func TestVoucherService_GetNextPage(t *testing.T) {
	service := &VoucherService{}

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

func TestVoucherlService_GetPrevPage(t *testing.T) {
	service := &VoucherService{}

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

func TestChallengeService_GetAll(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())

	service := NewVoucherService(repo, userService)

	vouchers := []*entities.VoucherModels{
		{
			Name:        "voucher a",
			ID:          1,
			Code:        "abc",
			Category:    "abc",
			Description: "abc",
			Discount:    1000,
			StartDate:   time.Now(),
			EndDate:     time.Now(),
			MinPurchase: 1000,
			Stock:       10,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
		{
			Name:        "voucher 2",
			ID:          2,
			Code:        "abc",
			Category:    "abc",
			Description: "abc",
			Discount:    1000,
			StartDate:   time.Now(),
			EndDate:     time.Now(),
			MinPurchase: 1000,
			Stock:       10,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
	}

	t.Run("Success Case - Voucher Found", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("FindAllVoucher", 1, 10).Return(vouchers, nil).Once()
		repo.On("GetTotalVoucherCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetAllVoucher(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(vouchers), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetTotalVoucherCount Error", func(t *testing.T) {
		expectedErr := errors.New("GetTotalVoucherCount Error")

		repo.On("FindAllVoucher", 1, 10).Return(vouchers, nil).Once()
		repo.On("GetTotalVoucherCount").Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetAllVoucher(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetAll Error", func(t *testing.T) {
		expectedErr := errors.New("FindAll Error")
		repo.On("FindAllVoucher", 1, 10).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetAllVoucher(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})
}

func TestVoucherService_Create(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repo, userService)

	existingvouchers := &entities.VoucherModels{
		Name:        "voucher a",
		ID:          1,
		Code:        "abc",
		Category:    "abc",
		Description: "abc",
		Discount:    1000,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		MinPurchase: 1000,
		Stock:       10,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	vouchers := &entities.VoucherModels{
		Name:        "voucher b",
		ID:          2,
		Code:        "abc",
		Category:    "abc",
		Description: "abc",
		Discount:    1000,
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		MinPurchase: 1000,
		Stock:       10,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	t.Run("Succes Case", func(t *testing.T) {
		repo.On("GetVoucherByCode", existingvouchers.Code).Return(nil, nil).Once()
		repo.On("CreateVoucher", mock.AnythingOfType("*entities.VoucherModels")).Return(nil, nil).Once()
		result, err := service.CreateVoucher(vouchers)
		assert.Nil(t, err)
		assert.Nil(t, result)
		assert.NotEqual(t, vouchers, result)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Creating Challenge", func(t *testing.T) {
		expectedErr := errors.New("failed to create challenge")
		repo.On("GetVoucherByCode", existingvouchers.Code).Return(nil, nil).Once()
		repo.On("CreateVoucher", mock.AnythingOfType("*entities.VoucherModels")).Return(nil, expectedErr).Once()
		result, err := service.CreateVoucher(vouchers)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.New("gagal menambahkan kupon"), err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - kupon sudah digunakan ", func(t *testing.T) {
		expectedErr := errors.New("failed to get voucher by code")
		repo.On("GetVoucherByCode", existingvouchers.Code).Return(existingvouchers, expectedErr).Once()
		result, err := service.CreateVoucher(vouchers)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.New("kode kupon sudah digunakan"), err)

		repo.AssertExpectations(t)
	})

	t.Run("Success Case - Create Voucher", func(t *testing.T) {
		existingVoucher := &entities.VoucherModels{
			Name:        "voucher a",
			ID:          1,
			Code:        "abc",
			Category:    "abc",
			Description: "abc",
			Discount:    1000,
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(time.Hour * 24),
			Status:      "Belum Kadaluwarsa",
		}

		newVoucher := &entities.VoucherModels{
			Name:        "voucher b",
			ID:          2,
			Code:        "def",
			Category:    "def",
			Description: "def",
			Discount:    1500,
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(time.Hour * 48),
		}

		repo.On("GetVoucherByCode", newVoucher.Code).Return(nil, nil).Once()
		repo.On("CreateVoucher", mock.AnythingOfType("*entities.VoucherModels")).Return(existingVoucher, nil).Once()

		result, err := service.CreateVoucher(newVoucher)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, existingVoucher, result)
		assert.Equal(t, "Belum Kadaluwarsa", result.Status)
		repo.AssertExpectations(t)
	})
}

func TestVoucherService_UpdateVoucher_Success(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repo, userService)

	existingVoucher := &entities.VoucherModels{
		Name:        "voucher a",
		ID:          1,
		Code:        "abc",
		Category:    "abc",
		Description: "abc",
		Discount:    1000,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * 24),
		MinPurchase: 1000,
		Stock:       10,
		Status:      "Belum Kadaluwarsa",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	updatedVoucher := &entities.VoucherModels{
		ID:          1,
		Name:        "voucher b",
		Code:        "def",
		Category:    "def",
		Description: "def",
		Discount:    1500,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * 48),
		MinPurchase: 1200,
		Stock:       15,
		Status:      "active",
		UpdatedAt:   time.Now(),
	}

	t.Run("Success Case - Update Voucher", func(t *testing.T) {
		repo.On("GetVoucherById", updatedVoucher.ID).Return(existingVoucher, nil).Once()
		repo.On("GetVoucherByCode", updatedVoucher.Code).Return(nil, nil).Once()
		repo.On("UpdateVoucher", updatedVoucher.ID, mock.AnythingOfType("*entities.VoucherModels")).Return(nil).Once()

		err := service.UpdateVoucher(updatedVoucher.ID, updatedVoucher)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Voucher Not Found", func(t *testing.T) {
		repo.On("GetVoucherById", updatedVoucher.ID).Return(nil, errors.New("kupon tidak ditemukan")).Once()

		err := service.UpdateVoucher(updatedVoucher.ID, updatedVoucher)

		assert.Error(t, err)
		assert.Equal(t, errors.New("kupon tidak ditemukan"), err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Voucher Code Already Used", func(t *testing.T) {
		repo.On("GetVoucherById", updatedVoucher.ID).Return(existingVoucher, nil).Once()
		repo.On("GetVoucherByCode", updatedVoucher.Code).Return(existingVoucher, nil).Once()

		err := service.UpdateVoucher(updatedVoucher.ID, updatedVoucher)

		assert.Error(t, err)
		assert.Equal(t, errors.New("kode kupon sudah digunakan"), err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Update Failure", func(t *testing.T) {
		repo.On("GetVoucherById", updatedVoucher.ID).Return(existingVoucher, nil).Once()
		repo.On("GetVoucherByCode", updatedVoucher.Code).Return(nil, nil).Once()
		repo.On("UpdateVoucher", updatedVoucher.ID, mock.AnythingOfType("*entities.VoucherModels")).Return(errors.New("gagal memperbarui kupon")).Once()

		err := service.UpdateVoucher(updatedVoucher.ID, updatedVoucher)

		assert.Error(t, err)
		assert.Equal(t, errors.New("gagal memperbarui kupon"), err)

		repo.AssertExpectations(t)
	})

	t.Run("Success Case - Update Voucher (Status Kadaluwarsa)", func(t *testing.T) {
		updatedVoucher := &entities.VoucherModels{
			ID:        1,
			Discount:  1500,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour * 48),
			Status:    "Kadaluwarsa",
			UpdatedAt: time.Now(),
		}

		updatedVoucher.EndDate = time.Now().Add(-time.Hour * 48)

		repo.On("GetVoucherById", updatedVoucher.ID).Return(existingVoucher, nil).Once()
		repo.On("GetVoucherByCode", updatedVoucher.Code).Return(nil, nil).Once()
		repo.On("UpdateVoucher", updatedVoucher.ID, mock.AnythingOfType("*entities.VoucherModels")).Return(nil).Once()

		err := service.UpdateVoucher(updatedVoucher.ID, updatedVoucher)

		assert.Nil(t, err)
		assert.Equal(t, "Kadaluwarsa", updatedVoucher.Status)

		repo.AssertExpectations(t)
	})
}

func TestVoucherService_DeleteVoucher(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repo, userService)

	t.Run("Success Delete", func(t *testing.T) {
		deletedVoucher := &entities.VoucherModels{
			ID:        1,
			Discount:  1500,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour * 48),
			Status:    "Kadaluwarsa",
			UpdatedAt: time.Now(),
		}
		repo.On("GetVoucherById", uint64(deletedVoucher.ID)).Return(deletedVoucher, nil).Once()
		repo.On("DeleteVoucher", uint64(deletedVoucher.ID)).Return(nil, nil).Once()
		err := service.DeleteVoucher(uint64(deletedVoucher.ID))

		assert.Nil(t, err)
	})

	t.Run("Failure Case - Voucher Not Found", func(t *testing.T) {
		// Specify a voucher ID that doesn't exist
		nonExistentVoucherID := uint64(999)

		// Mock the GetVoucherById function to return an error for the non-existent voucher
		repo.On("GetVoucherById", nonExistentVoucherID).Return(nil, errors.New("kupon tidak ditemukan")).Once()

		err := service.DeleteVoucher(nonExistentVoucherID)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "kupon tidak ditemukan")
	})

	t.Run("Failure Case - Error Deleting Voucher", func(t *testing.T) {
		// Specify a voucher ID for which deletion will fail
		failingVoucherID := uint64(888)

		// Mock the GetVoucherById function to return a voucher
		voucher := &entities.VoucherModels{ID: failingVoucherID}
		repo.On("GetVoucherById", failingVoucherID).Return(voucher, nil).Once()

		// Mock the DeleteVoucher function to return an error
		repo.On("DeleteVoucher", failingVoucherID).Return(errors.New("error deleting voucher")).Once()

		err := service.DeleteVoucher(failingVoucherID)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error deleting voucher")
	})

}

func TestService_GetVoucherById(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repo, userService)

	t.Run("Success - Voucher Found", func(t *testing.T) {

		// Mocked voucher data
		mockVoucher := &entities.VoucherModels{
			ID:        1,
			Discount:  1500,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour * 48),
			Status:    "Active",
			UpdatedAt: time.Now(),
		}

		// Set up the repository mock to return the mocked voucher
		repo.On("GetVoucherById", mockVoucher.ID).Return(mockVoucher, nil).Once()

		// Call the GetVoucherById function
		result, err := service.GetVoucherById(mockVoucher.ID)

		// Assert that there is no error
		assert.Nil(t, err)

		// Assert that the returned voucher matches the mocked voucher
		assert.Equal(t, mockVoucher, result)

	})

	t.Run("Failure Case - Voucher Not Found", func(t *testing.T) {
		// Specify a voucher ID that doesn't exist
		nonExistentVoucherID := &entities.VoucherModels{
			ID:        1,
			Discount:  1500,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour * 48),
			Status:    "Active",
			UpdatedAt: time.Now(),
		}

		// Mock the GetVoucherById function to return an error for the non-existent voucher
		repo.On("GetVoucherById", uint64(nonExistentVoucherID.ID)).Return(nil, errors.New("kupon tidak ditemukan")).Once()

		_, err := service.GetVoucherById(uint64(nonExistentVoucherID.ID))

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)

		assert.EqualError(t, err, "kupon tidak ditemukan")
	})
}

func TestVoucher_DeleteVoucherClaims(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repo, userService)

	t.Run("Success - Delete User Voucher Claims", func(t *testing.T) {
		// Mocked user and voucher IDs
		userID := uint64(123)
		voucherID := uint64(456)

		// Create an instance of VoucherService with mocked dependencies

		repo.On("DeleteUserVoucherClaims", userID, voucherID).Return(nil).Once()

		// Call the DeleteUserVoucherClaims function
		err := service.DeleteVoucherClaims(userID, voucherID)

		// Assert that there is no error
		assert.Nil(t, err)

		// Assert that the mocked method was called as expected
		repo.AssertExpectations(t)
	})

	t.Run("Failure - Error Deleting User Voucher Claims", func(t *testing.T) {
		// Mocked user and voucher IDs
		userID := uint64(123)
		voucherID := uint64(456)

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks to simulate an error when deleting user voucher claims

		repo.On("DeleteUserVoucherClaims", userID, voucherID).Return(errors.New("error deleting user voucher claims")).Once()

		// Call the DeleteUserVoucherClaims function
		err := service.DeleteVoucherClaims(userID, voucherID)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error deleting user voucher claims")

		// Assert that the other methods were not called
		repo.AssertExpectations(t)
	})
}

func TestGetUserVouchers(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repo, userService)
	t.Run("Success - Get User Vouchers", func(t *testing.T) {
		// Mocked user ID
		userID := uint64(123)

		// Mocked user vouchers
		mockUserVouchers := []*entities.VoucherClaimModels{
			{UserID: userID, VoucherID: 1},
			{UserID: userID, VoucherID: 2},
		}

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks

		repo.On("GetUserVoucherClaims", userID).Return(mockUserVouchers, nil).Once()

		// Call the GetUserVouchers function
		userVouchers, err := service.GetUserVouchers(userID)

		// Assert that there is no error
		assert.Nil(t, err)

		// Assert that the returned user vouchers match the mocked user vouchers
		assert.Equal(t, mockUserVouchers, userVouchers)

		// Assert that the mocked method was called as expected
		repo.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting User Vouchers", func(t *testing.T) {
		// Mocked user ID
		userID := uint64(123)

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks to simulate an error when getting user vouchers

		repo.On("GetUserVoucherClaims", userID).Return(nil, errors.New("error getting user vouchers")).Once()

		// Call the GetUserVouchers function
		userVouchers, err := service.GetUserVouchers(userID)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting user vouchers")

		// Assert that the result is nil when there is an error
		assert.Nil(t, userVouchers)

		// Assert that the other methods were not called
		repo.AssertExpectations(t)
	})
}

func TestGetVoucherByStatus(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repoMock, userService)
	t.Run("Success - Get Voucher By Status", func(t *testing.T) {
		// Mocked parameters
		page := 1
		perPage := 10
		status := "Active"

		// Mocked vouchers and total items
		mockVouchers := []*entities.VoucherModels{
			{ID: 1, Status: "Active"},
			{ID: 2, Status: "Active"},
		}
		mockTotalItems := int64(2)

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks

		repoMock.On("FindByStatus", page, perPage, status).Return(mockVouchers, nil).Once()
		repoMock.On("GetTotalVoucherCountByStatus", status).Return(mockTotalItems, nil).Once()

		// Call the GetVoucherByStatus function
		vouchers, totalItems, err := service.GetVoucherByStatus(page, perPage, status)

		// Assert that there is no error
		assert.Nil(t, err)

		// Assert that the returned vouchers and total items match the mocked values
		assert.Equal(t, mockVouchers, vouchers)
		assert.Equal(t, mockTotalItems, totalItems)

		// Assert that the mocked methods were called as expected
		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Vouchers by Status", func(t *testing.T) {
		// Mocked parameters
		page := 1
		perPage := 10
		status := "Expired"

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks to simulate an error when getting vouchers by status

		repoMock.On("FindByStatus", page, perPage, status).Return(nil, errors.New("error getting vouchers by status")).Once()

		// Call the GetVoucherByStatus function
		vouchers, totalItems, err := service.GetVoucherByStatus(page, perPage, status)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers by status")

		// Assert that the result is nil when there is an error
		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		// Assert that the other methods were not called
		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Total Voucher Count", func(t *testing.T) {
		// Mocked parameters
		page := 1
		perPage := 10
		status := "active"

		// Set up mocks to simulate an error when getting the total voucher count
		repoMock.On("FindByStatus", page, perPage, status).Return(nil, errors.New("error getting vouchers")).Once()
		repoMock.On("GetTotalVoucherCountByStatus", status).Return(int64(0), errors.New("error getting total voucher count")).Once()

		// Call the GetVoucherByStatus function
		vouchers, totalItems, err := service.GetVoucherByStatus(page, perPage, status)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers")

		// Assert that the result is nil when there is an error
		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		// Assert that the other methods were not called
		repoMock.AssertExpectations(t)
	})
}

func TestGetVoucherByCategory(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repoMock, userService)
	t.Run("Success - Get Voucher By Category", func(t *testing.T) {
		// Mocked parameters
		page := 1
		perPage := 10
		category := "Discount"

		// Mocked vouchers and total items
		mockVouchers := []*entities.VoucherModels{
			{ID: 1, Category: "Discount"},
			{ID: 2, Category: "Discount"},
		}
		mockTotalItems := int64(2)

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks

		repoMock.On("FindByCategory", page, perPage, category).Return(mockVouchers, nil).Once()
		repoMock.On("GetTotalVoucherCountByCategory", category).Return(mockTotalItems, nil).Once()

		// Call the GetVoucherByCategory function
		vouchers, totalItems, err := service.GetVoucherByCategory(page, perPage, category)

		// Assert that there is no error
		assert.Nil(t, err)

		// Assert that the returned vouchers and total items match the mocked values
		assert.Equal(t, mockVouchers, vouchers)
		assert.Equal(t, mockTotalItems, totalItems)

		// Assert that the mocked methods were called as expected
		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Vouchers by Category", func(t *testing.T) {
		// Mocked parameters
		page := 1
		perPage := 10
		category := "NonexistentCategory"

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks to simulate an error when getting vouchers by category

		repoMock.On("FindByCategory", page, perPage, category).Return(nil, errors.New("error getting vouchers by category")).Once()

		// Call the GetVoucherByCategory function
		vouchers, totalItems, err := service.GetVoucherByCategory(page, perPage, category)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers by category")

		// Assert that the result is nil when there is an error
		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		// Assert that the other methods were not called
		repoMock.AssertExpectations(t)
	})

}

func TestGetVoucherByStatusCategory(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repoMock, userService)
	t.Run("Success - Get Voucher By Status and Category", func(t *testing.T) {
		// Mocked parameters
		page := 1
		perPage := 10
		status := "Active"
		category := "Discount"

		// Mocked vouchers and total items
		mockVouchers := []*entities.VoucherModels{
			{ID: 1, Status: "Active", Category: "Discount"},
			{ID: 2, Status: "Active", Category: "Discount"},
		}
		mockTotalItems := int64(2)

		// Create an instance of VoucherService with mocked dependencies

		repoMock.On("FindByStatusCategory", page, perPage, status, category).Return(mockVouchers, nil).Once()
		repoMock.On("GetTotalVoucherCountByStatusCategory", status, category).Return(mockTotalItems, nil).Once()

		// Call the GetVoucherByStatusCategory function
		vouchers, totalItems, err := service.GetVoucherByStatusCategory(page, perPage, status, category)

		// Assert that there is no error
		assert.Nil(t, err)

		// Assert that the returned vouchers and total items match the mocked values
		assert.Equal(t, mockVouchers, vouchers)
		assert.Equal(t, mockTotalItems, totalItems)

		// Assert that the mocked methods were called as expected
		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Vouchers by Status and Category", func(t *testing.T) {
		// Mocked parameters
		page := 1
		perPage := 10
		status := "Expired"
		category := "NonexistentCategory"

		// Set up mocks to simulate an error when getting vouchers by status and category

		repoMock.On("FindByStatusCategory", page, perPage, status, category).Return(nil, errors.New("error getting vouchers by status and category")).Once()

		// Call the GetVoucherByStatusCategory function
		vouchers, totalItems, err := service.GetVoucherByStatusCategory(page, perPage, status, category)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers by status and category")

		// Assert that the result is nil when there is an error
		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		// Assert that the other methods were not called
		repoMock.AssertExpectations(t)
	})
}

func TestGetAllVoucherToClaims(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repoMock, userService)
	t.Run("Success - Get All Vouchers to Claims", func(t *testing.T) {
		// Mocked parameters
		limit := 5
		userID := uint64(123)

		// Mocked vouchers
		mockVouchers := []*entities.VoucherModels{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		}

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks

		repoMock.On("FindAllVoucherToClaims", limit, userID).Return(mockVouchers, nil).Once()
		repoMock.On("IsVoucherAlreadyClaimed", userID, mock.AnythingOfType("uint64")).Return(false, nil).Times(len(mockVouchers))

		// Call the GetAllVoucherToClaims function
		filteredVouchers, err := service.GetAllVoucherToClaims(limit, userID)

		// Assert that there is no error
		assert.Nil(t, err)

		// Assert that the number of filtered vouchers matches the expected count
		assert.Len(t, filteredVouchers, len(mockVouchers))

		// Assert that each voucher is not claimed
		for _, voucher := range filteredVouchers {
			assert.False(t, voucher.Status == "active")
		}

		// Assert that the mocked methods were called as expected
		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Vouchers to Claims", func(t *testing.T) {
		// Mocked parameters
		limit := 5
		userID := uint64(123)

		// Create an instance of VoucherService with mocked dependencies

		// Set up mocks to simulate an error when getting vouchers to claims

		repoMock.On("FindAllVoucherToClaims", limit, userID).Return(nil, errors.New("error getting vouchers to claims")).Once()

		// Call the GetAllVoucherToClaims function
		filteredVouchers, err := service.GetAllVoucherToClaims(limit, userID)

		// Assert that the error is not nil and it matches the expected error message
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers to claims")

		// Assert that the result is nil when there is an error
		assert.Nil(t, filteredVouchers)

		// Assert that the other methods were not called
		repoMock.AssertExpectations(t)
	})
}
