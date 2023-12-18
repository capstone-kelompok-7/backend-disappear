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

func TestVoucherService_GetAll(t *testing.T) {
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

	t.Run("Failure Case - Claimed Check Error", func(t *testing.T) {
		limit := 10
		userID := uint64(1)
		expectedErr := errors.New("IsVoucherAlreadyClaimed Error")

		repo.On("FindAllVoucherToClaims", limit, userID).Return(vouchers, nil).Once()
		repo.On("IsVoucherAlreadyClaimed", userID, mock.Anything).Return(false, expectedErr).Once()

		result, err := service.GetAllVoucherToClaims(limit, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "IsVoucherAlreadyClaimed Error")
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
		nonExistentVoucherID := uint64(999)

		repo.On("GetVoucherById", nonExistentVoucherID).Return(nil, errors.New("kupon tidak ditemukan")).Once()

		err := service.DeleteVoucher(nonExistentVoucherID)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "kupon tidak ditemukan")
	})

	t.Run("Failure Case - Error Deleting Voucher", func(t *testing.T) {
		failingVoucherID := uint64(888)

		voucher := &entities.VoucherModels{ID: failingVoucherID}
		repo.On("GetVoucherById", failingVoucherID).Return(voucher, nil).Once()
		repo.On("DeleteVoucher", failingVoucherID).Return(errors.New("error deleting voucher")).Once()

		err := service.DeleteVoucher(failingVoucherID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error deleting voucher")
	})

}

func TestVoucher_GetVoucherById(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repo, userService)

	t.Run("Success - Voucher Found", func(t *testing.T) {

		mockVoucher := &entities.VoucherModels{
			ID:        1,
			Discount:  1500,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour * 48),
			Status:    "Active",
			UpdatedAt: time.Now(),
		}

		repo.On("GetVoucherById", mockVoucher.ID).Return(mockVoucher, nil).Once()

		result, err := service.GetVoucherById(mockVoucher.ID)

		assert.Nil(t, err)
		assert.Equal(t, mockVoucher, result)
	})

	t.Run("Failure Case - Voucher Not Found", func(t *testing.T) {
		nonExistentVoucherID := &entities.VoucherModels{
			ID:        1,
			Discount:  1500,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour * 48),
			Status:    "Active",
			UpdatedAt: time.Now(),
		}

		repo.On("GetVoucherById", uint64(nonExistentVoucherID.ID)).Return(nil, errors.New("kupon tidak ditemukan")).Once()

		_, err := service.GetVoucherById(uint64(nonExistentVoucherID.ID))

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

func TestVoucher_TestGetUserVouchers(t *testing.T) {
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

func TestVoucher_TestGetVoucherByStatus(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repoMock, userService)

	t.Run("Success - Get Voucher By Status", func(t *testing.T) {
		page := 1
		perPage := 10
		status := "Active"

		mockVouchers := []*entities.VoucherModels{
			{ID: 1, Status: "Active"},
			{ID: 2, Status: "Active"},
		}
		mockTotalItems := int64(2)

		repoMock.On("FindByStatus", page, perPage, status).Return(mockVouchers, nil).Once()
		repoMock.On("GetTotalVoucherCountByStatus", status).Return(mockTotalItems, nil).Once()
		vouchers, totalItems, err := service.GetVoucherByStatus(page, perPage, status)

		assert.Nil(t, err)
		assert.Equal(t, mockVouchers, vouchers)
		assert.Equal(t, mockTotalItems, totalItems)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Vouchers by Status", func(t *testing.T) {
		page := 1
		perPage := 10
		status := "Expired"

		repoMock.On("FindByStatus", page, perPage, status).Return(nil, errors.New("error getting vouchers by status")).Once()
		vouchers, totalItems, err := service.GetVoucherByStatus(page, perPage, status)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers by status")
		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Total Voucher Count by Status", func(t *testing.T) {
		page := 1
		perPage := 10
		status := "Active"

		mockVouchers := []*entities.VoucherModels{
			{ID: 1, Status: "Active"},
			{ID: 2, Status: "Active"},
		}

		repoMock.On("FindByStatus", page, perPage, status).Return(mockVouchers, nil).Once()
		repoMock.On("GetTotalVoucherCountByStatus", status).Return(int64(0), errors.New("error getting total voucher count")).Once()

		vouchers, totalItems, err := service.GetVoucherByStatus(page, perPage, status)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting total voucher count")
		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		repoMock.AssertExpectations(t)
	})
}

func TestVoucher_TestGetVoucherByCategory(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repoMock, userService)

	page := 1
	perPage := 10
	category := "Discount"

	t.Run("Success - Get Voucher By Category", func(t *testing.T) {
		mockVouchers := []*entities.VoucherModels{
			{ID: 1, Category: "Discount"},
			{ID: 2, Category: "Discount"},
		}
		mockTotalItems := int64(2)

		repoMock.On("FindByCategory", page, perPage, category).Return(mockVouchers, nil).Once()
		repoMock.On("GetTotalVoucherCountByCategory", category).Return(mockTotalItems, nil).Once()

		vouchers, totalItems, err := service.GetVoucherByCategory(page, perPage, category)

		assert.Nil(t, err)
		assert.Equal(t, mockVouchers, vouchers)
		assert.Equal(t, mockTotalItems, totalItems)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Vouchers by Category", func(t *testing.T) {
		page := 1
		perPage := 10
		category := "NonexistentCategory"

		repoMock.On("FindByCategory", page, perPage, category).Return(nil, errors.New("error getting vouchers by category")).Once()

		vouchers, totalItems, err := service.GetVoucherByCategory(page, perPage, category)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers by category")

		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Total Voucher Count by Category", func(t *testing.T) {
		expectedError := errors.New("error getting total voucher count by category")
		repoMock.On("FindByCategory", page, perPage, category).Return(nil, nil).Once()
		repoMock.On("GetTotalVoucherCountByCategory", category).Return(int64(0), expectedError).Once()

		vouchers, totalItems, err := service.GetVoucherByCategory(page, perPage, category)

		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedError.Error())

		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		repoMock.AssertExpectations(t)
	})
}

func TestVoucher_TestGetVoucherByStatusCategory(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repoMock, userService)

	page := 1
	perPage := 10
	status := "Active"
	category := "Discount"

	t.Run("Success - Get Voucher By Status and Category", func(t *testing.T) {
		mockVouchers := []*entities.VoucherModels{
			{ID: 1, Status: "Active", Category: "Discount"},
			{ID: 2, Status: "Active", Category: "Discount"},
		}
		mockTotalItems := int64(2)

		repoMock.On("FindByStatusCategory", page, perPage, status, category).Return(mockVouchers, nil).Once()
		repoMock.On("GetTotalVoucherCountByStatusCategory", status, category).Return(mockTotalItems, nil).Once()

		vouchers, totalItems, err := service.GetVoucherByStatusCategory(page, perPage, status, category)

		assert.Nil(t, err)
		assert.Equal(t, mockVouchers, vouchers)
		assert.Equal(t, mockTotalItems, totalItems)
		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Vouchers by Status and Category", func(t *testing.T) {
		page := 1
		perPage := 10
		status := "Expired"
		category := "NonexistentCategory"

		repoMock.On("FindByStatusCategory", page, perPage, status, category).Return(nil, errors.New("error getting vouchers by status and category")).Once()

		vouchers, totalItems, err := service.GetVoucherByStatusCategory(page, perPage, status, category)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers by status and category")
		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)
		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Total Voucher Count by Status and Category", func(t *testing.T) {
		expectedError := errors.New("error getting total voucher count by status and category")
		repoMock.On("FindByStatusCategory", page, perPage, status, category).Return(nil, nil).Once()
		repoMock.On("GetTotalVoucherCountByStatusCategory", status, category).Return(int64(0), expectedError).Once()

		vouchers, totalItems, err := service.GetVoucherByStatusCategory(page, perPage, status, category)

		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedError.Error())

		assert.Nil(t, vouchers)
		assert.Zero(t, totalItems)

		repoMock.AssertExpectations(t)
	})
}

func TestVoucher_TestGetAllVoucherToClaims(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewVoucherService(repoMock, userService)
	t.Run("Success - Get All Vouchers to Claims", func(t *testing.T) {
		limit := 5
		userID := uint64(123)

		mockVouchers := []*entities.VoucherModels{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		}

		repoMock.On("FindAllVoucherToClaims", limit, userID).Return(mockVouchers, nil).Once()
		repoMock.On("IsVoucherAlreadyClaimed", userID, mock.AnythingOfType("uint64")).Return(false, nil).Times(len(mockVouchers))

		filteredVouchers, err := service.GetAllVoucherToClaims(limit, userID)

		assert.Nil(t, err)

		assert.Len(t, filteredVouchers, len(mockVouchers))

		for _, voucher := range filteredVouchers {
			assert.False(t, voucher.Status == "active")
		}

		repoMock.AssertExpectations(t)
	})

	t.Run("Failure - Error Getting Vouchers to Claims", func(t *testing.T) {
		limit := 5
		userID := uint64(123)

		repoMock.On("FindAllVoucherToClaims", limit, userID).Return(nil, errors.New("error getting vouchers to claims")).Once()

		filteredVouchers, err := service.GetAllVoucherToClaims(limit, userID)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting vouchers to claims")
		assert.Nil(t, filteredVouchers)
		repoMock.AssertExpectations(t)
	})
}

func TestVoucher_TestCanClaimsVoucher(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	userMock := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(userMock, utils.NewHash())
	service := NewVoucherService(repoMock, userService)

	t.Run("Success - User can claim voucher", func(t *testing.T) {
		userID := uint64(1)
		voucherID := uint64(100)
		userLevel := "Gold"

		mockUser := &entities.UserModels{ID: userID, Level: userLevel}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return(userLevel, nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Gold", nil).Once()

		canClaim, err := service.CanClaimsVoucher(userID, voucherID)

		assert.Nil(t, err)
		assert.True(t, canClaim)

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - Error getting user level", func(t *testing.T) {
		userID := uint64(1)
		voucherID := uint64(100)
		mockUser := &entities.UserModels{ID: userID}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return("", errors.New("error getting user level")).Once()

		canClaim, err := service.CanClaimsVoucher(userID, voucherID)

		assert.NotNil(t, err)
		assert.False(t, canClaim)
		assert.EqualError(t, err, "error getting user level")

		userMock.AssertExpectations(t)
	})

	t.Run("Failure - Error getting voucher category", func(t *testing.T) {
		userID := uint64(1)
		voucherID := uint64(100)
		userLevel := "Gold"

		mockUser := &entities.UserModels{ID: userID, Level: userLevel}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return(userLevel, nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("", errors.New("error getting voucher category")).Once()

		canClaim, err := service.CanClaimsVoucher(userID, voucherID)

		assert.NotNil(t, err)
		assert.False(t, canClaim)
		assert.EqualError(t, err, "error getting voucher category")

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - User cannot claim voucher", func(t *testing.T) {
		userID := uint64(1)
		voucherID := uint64(100)
		userLevel := "Silver"

		mockUser := &entities.UserModels{ID: userID, Level: userLevel}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return(userLevel, nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Gold", nil).Once()

		canClaim, err := service.CanClaimsVoucher(userID, voucherID)

		assert.Nil(t, err)
		assert.False(t, canClaim)

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - Voucher Category is All Customer", func(t *testing.T) {
		userID := uint64(1)
		voucherID := uint64(100)

		mockUser := &entities.UserModels{ID: userID, Level: "Gold"}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return("Gold", nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("All Customer", nil).Once()

		canClaim, err := service.CanClaimsVoucher(userID, voucherID)

		assert.Nil(t, err)
		assert.True(t, canClaim, "user should not be able to claim voucher with category All Customer")

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Success - User Level is Silver and Voucher Category is Bronze", func(t *testing.T) {
		userID := uint64(1)
		voucherID := uint64(100)

		mockUser := &entities.UserModels{ID: userID, Level: "Silver"}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return("Silver", nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Bronze", nil).Once()

		canClaim, err := service.CanClaimsVoucher(userID, voucherID)

		assert.Nil(t, err)
		assert.True(t, canClaim)

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Success - User Level is Bronze and Voucher Category is Bronze", func(t *testing.T) {
		userID := uint64(1)
		voucherID := uint64(100)

		mockUser := &entities.UserModels{ID: userID, Level: "Bronze"}

		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return("Bronze", nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Bronze", nil).Once()

		canClaim, err := service.CanClaimsVoucher(userID, voucherID)

		assert.Nil(t, err)
		assert.True(t, canClaim)

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})
}

func TestVoucher_TestClaimVoucher(t *testing.T) {
	repoMock := voucherMocks.NewRepositoryVoucherInterface(t)
	userMock := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(userMock, utils.NewHash())
	service := NewVoucherService(repoMock, userService)

	userID := uint64(1)
	voucherID := uint64(100)

	mockVoucher := &entities.VoucherModels{
		ID:    voucherID,
		Stock: uint64(100),
	}

	t.Run("Success - Claim Voucher", func(t *testing.T) {

		mockUser := &entities.UserModels{ID: userID, Level: "Gold"}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		repoMock.On("GetVoucherById", voucherID).Return(mockVoucher, nil).Once()
		repoMock.On("IsVoucherAlreadyClaimed", userID, voucherID).Return(false, nil).Once()
		repoMock.On("ClaimVoucher", &entities.VoucherClaimModels{UserID: userID, VoucherID: voucherID}).Return(nil).Once()
		repoMock.On("ReduceStockWhenClaimed", voucherID, uint64(1)).Return(nil).Once()

		userMock.On("GetUserLevel", userID).Return("Gold", nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Gold", nil).Once()

		req := &entities.VoucherClaimModels{
			UserID:    userID,
			VoucherID: voucherID,
		}

		err := service.ClaimVoucher(req)

		assert.Nil(t, err)

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - Out of Stock", func(t *testing.T) {
		mockVoucher := &entities.VoucherModels{
			ID:    voucherID,
			Stock: uint64(0),
		}
		repoMock.On("GetVoucherById", voucherID).Return(mockVoucher, nil).Once()

		req := &entities.VoucherClaimModels{
			UserID:    userID,
			VoucherID: voucherID,
		}

		err := service.ClaimVoucher(req)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "stok kupon sudah habis")

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - Voucher Not Found", func(t *testing.T) {
		repoMock.On("GetVoucherById", voucherID).Return(nil, errors.New("kupon tidak ditemukan")).Once()

		req := &entities.VoucherClaimModels{
			UserID:    userID,
			VoucherID: voucherID,
		}

		err := service.ClaimVoucher(req)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "kupon tidak ditemukan")

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - IsVoucherAlreadyClaimed Error", func(t *testing.T) {
		mockUser := &entities.UserModels{ID: userID, Level: "Gold"}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return("Gold", nil).Once()
		repoMock.On("GetVoucherById", voucherID).Return(mockVoucher, nil).Once()
		repoMock.On("IsVoucherAlreadyClaimed", userID, voucherID).Return(false, errors.New("IsVoucherAlreadyClaimed error")).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Gold", nil).Once()
		req := &entities.VoucherClaimModels{
			UserID:    userID,
			VoucherID: voucherID,
		}

		err := service.ClaimVoucher(req)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "IsVoucherAlreadyClaimed error")
		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - Voucher Already Claimed", func(t *testing.T) {
		mockUser := &entities.UserModels{ID: userID, Level: "Gold"}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return("Gold", nil).Once()
		repoMock.On("GetVoucherById", voucherID).Return(mockVoucher, nil).Once()
		repoMock.On("IsVoucherAlreadyClaimed", userID, voucherID).Return(true, nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Gold", nil).Once()
		req := &entities.VoucherClaimModels{
			UserID:    userID,
			VoucherID: voucherID,
		}

		err := service.ClaimVoucher(req)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "kupon telah diklaim")

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - Insufficient Level", func(t *testing.T) {
		mockUser := &entities.UserModels{ID: userID, Level: "Bronze"}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		userMock.On("GetUserLevel", userID).Return("Bronze", nil).Once()
		repoMock.On("GetVoucherById", voucherID).Return(mockVoucher, nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Gold", nil).Once()

		req := &entities.VoucherClaimModels{
			UserID:    userID,
			VoucherID: voucherID,
		}

		err := service.ClaimVoucher(req)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "level anda masih belum mencukupi")

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})

	t.Run("Failure - ClaimVoucher Error", func(t *testing.T) {
		expectedError := errors.New("Expected ClaimVoucher error")
		mockUser := &entities.UserModels{ID: userID, Level: "Gold"}
		userMock.On("GetUsersById", mock.Anything).Return(mockUser, nil).Once()
		repoMock.On("GetVoucherById", voucherID).Return(mockVoucher, nil).Once()
		repoMock.On("IsVoucherAlreadyClaimed", userID, voucherID).Return(false, nil).Once()
		repoMock.On("ClaimVoucher", mock.Anything).Return(expectedError).Once()
		userMock.On("GetUserLevel", userID).Return("Gold", nil).Once()
		repoMock.On("GetVoucherCategory", voucherID).Return("Gold", nil).Once()

		req := &entities.VoucherClaimModels{
			UserID:    userID,
			VoucherID: voucherID,
		}

		err := service.ClaimVoucher(req)

		assert.EqualError(t, err, expectedError.Error())

		repoMock.AssertExpectations(t)
		userMock.AssertExpectations(t)
	})
}
