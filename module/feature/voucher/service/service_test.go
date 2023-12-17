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

func TestVoucherService_GetChallengeByTitle(t *testing.T) {
	repo := voucherMocks.NewRepositoryVoucherInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())

	service := NewVoucherService(repo, userService)

	vouchers := &entities.VoucherModels{
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

	t.Run("Success Case - Challenge Found by Name", func(t *testing.T) {
		repo.On("GetVoucherById", uint64(1)).Return(vouchers, nil).Once()

		resultFound, errFound := service.GetVoucherById(uint64(1))

		assert.Nil(t, errFound)
		assert.NotNil(t, resultFound)
		assert.Equal(t, vouchers, resultFound)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Challenge by Title", func(t *testing.T) {
		expectedErr := errors.New("failed to find challenge by tittle")
		repo.On("GetVoucherById", uint64(1)).Return(nil, expectedErr).Once()

		result, errNotFound := service.GetVoucherById(uint64(1))

		assert.Error(t, errNotFound)
		assert.Nil(t, result)
		assert.Equal(t, errors.New("kupon tidak ditemukan"), errNotFound)

		repo.AssertExpectations(t)
	})
}

func TestVoucherService_CreateSubmitChallengeForm(t *testing.T) {
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
		assert.NotNil(t, result)
		assert.Equal(t, vouchers, result)

		repo.AssertExpectations(t)
	})

}
