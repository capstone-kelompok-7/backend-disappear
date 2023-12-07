package service

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard/mocks"
	mocks_caching "github.com/capstone-kelompok-7/backend-disappear/utils/caching/mocks"
	"github.com/stretchr/testify/assert"
)

func setupTestService(t *testing.T) (*mocks.RepositoryDashboardInterface, *mocks_caching.CacheRepository, dashboard.ServiceDashboardInterface) {
	repo := mocks.NewRepositoryDashboardInterface(t)
	cache := mocks_caching.NewCacheRepository(t)
	service := NewDashboardService(repo, cache)
	return repo, cache, service
}

func TestGetCardDashboard(t *testing.T) {
	expectedProductCount := int64(10)
	expectedOrderCount := int64(20)
	expectedUserCount := int64(30)
	expectedIncomeCount := float64(5000.0)

	t.Run("Success Case", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		repo.On("CountProducts").Return(expectedProductCount, nil)
		repo.On("CountOrder").Return(expectedOrderCount, nil)
		repo.On("CountUsers").Return(expectedUserCount, nil)
		repo.On("CountIncome").Return(expectedIncomeCount, nil)

		productCount, userCount, orderCount, incomeCount, err := service.GetCardDashboard()

		assert.NoError(t, err)
		assert.Equal(t, expectedProductCount, productCount)
		assert.Equal(t, expectedUserCount, userCount)
		assert.Equal(t, expectedOrderCount, orderCount)
		assert.Equal(t, expectedIncomeCount, incomeCount)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CountProducts", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal menghitung total produk")
		repo.On("CountProducts").Return(int64(0), expectedError)

		productCount, userCount, orderCount, incomeCount, err := service.GetCardDashboard()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, productCount)
		assert.Zero(t, userCount)
		assert.Zero(t, orderCount)
		assert.Zero(t, incomeCount)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CountOrder", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal menghitung total pesanan")
		repo.On("CountProducts").Return(int64(10), nil)
		repo.On("CountOrder").Return(int64(0), expectedError)

		productCount, userCount, orderCount, incomeCount, err := service.GetCardDashboard()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, productCount)
		assert.Zero(t, userCount)
		assert.Zero(t, orderCount)
		assert.Zero(t, incomeCount)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CountUsers", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal menghitung total pelanggan")
		repo.On("CountProducts").Return(int64(10), nil)
		repo.On("CountOrder").Return(int64(20), nil)
		repo.On("CountUsers").Return(int64(0), expectedError)

		productCount, userCount, orderCount, incomeCount, err := service.GetCardDashboard()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, productCount)
		assert.Zero(t, userCount)
		assert.Zero(t, orderCount)
		assert.Zero(t, incomeCount)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CountIncome", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal menghitung total pendapatan")
		repo.On("CountProducts").Return(int64(10), nil)
		repo.On("CountOrder").Return(int64(20), nil)
		repo.On("CountUsers").Return(int64(30), nil)
		repo.On("CountIncome").Return(float64(0.0), expectedError)

		productCount, userCount, orderCount, incomeCount, err := service.GetCardDashboard()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, productCount)
		assert.Zero(t, userCount)
		assert.Zero(t, orderCount)
		assert.Zero(t, incomeCount)
		repo.AssertExpectations(t)
	})
}

func TestGetLandingPage(t *testing.T) {
	expectedUserCount := int64(30)
	expectedGramPlastic := int64(100)
	expectedOrderCount := int64(20)

	t.Run("Success Case", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		repo.On("CountUsers").Return(expectedUserCount, nil)
		repo.On("CountTotalGram").Return(expectedGramPlastic, nil)
		repo.On("CountOrder").Return(expectedOrderCount, nil)

		userCount, gramPlastic, orderCount, err := service.GetLandingPage()

		assert.NoError(t, err)
		assert.Equal(t, expectedUserCount, userCount)
		assert.Equal(t, expectedGramPlastic, gramPlastic)
		assert.Equal(t, expectedOrderCount, orderCount)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CountUsers", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal menghitung total pelanggan")
		repo.On("CountUsers").Return(int64(0), expectedError)

		userCount, gramPlastic, orderCount, err := service.GetLandingPage()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, userCount)
		assert.Zero(t, gramPlastic)
		assert.Zero(t, orderCount)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CountTotalGram", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal menghitung total gram plastik")
		repo.On("CountUsers").Return(expectedUserCount, nil)
		repo.On("CountTotalGram").Return(int64(0), expectedError)

		userCount, gramPlastic, orderCount, err := service.GetLandingPage()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, userCount)
		assert.Zero(t, gramPlastic)
		assert.Zero(t, orderCount)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CountOrder", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal menghitung total pesanan")
		repo.On("CountUsers").Return(expectedUserCount, nil)
		repo.On("CountTotalGram").Return(expectedGramPlastic, nil)
		repo.On("CountOrder").Return(int64(0), expectedError)

		userCount, gramPlastic, orderCount, err := service.GetLandingPage()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, userCount)
		assert.Zero(t, gramPlastic)
		assert.Zero(t, orderCount)

		repo.AssertExpectations(t)
	})
}

func TestGetProductReviewsWithMaxTotal(t *testing.T) {
	mockProducts := []*entities.ProductModels{
		{
			ID:          1,
			Name:        "Product A",
			Description: "Description A",
			TotalReview: 10,
		},
	}

	t.Run("Success Case", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		repo.On("GetProductWithMaxReviews").Return(mockProducts, nil)
		response, err := service.GetProductReviewsWithMaxTotal()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, len(mockProducts), len(response))
		assert.Equal(t, response[0].ID, uint64(1))
		assert.Equal(t, response[0].Name, "Product A")
		assert.Equal(t, response[0].Description, "Description A")
		assert.Equal(t, response[0].TotalReview, uint64(10))
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal mengambil data produk dengan review terbanyak")
		repo.On("GetProductWithMaxReviews").Return(nil, expectedError)

		response, err := service.GetProductReviewsWithMaxTotal()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, expectedError.Error())

		repo.AssertExpectations(t)
	})
}

func TestGetGramPlasticStat(t *testing.T) {
	startOfWeek := time.Now().AddDate(0, 0, -7)
	endOfWeek := time.Now()

	expectedGramTotalCount := uint64(100)

	t.Run("Success Case - Cache Hit", func(t *testing.T) {
		_, cache, service := setupTestService(t)
		cacheKey := generateCacheKey(startOfWeek, endOfWeek)
		cache.On("Get", cacheKey).Return([]byte(strconv.FormatUint(expectedGramTotalCount, 10)), nil)

		gramTotalCount, err := service.GetGramPlasticStat(startOfWeek, endOfWeek)

		assert.NoError(t, err)
		assert.Equal(t, expectedGramTotalCount, gramTotalCount)

		cache.AssertExpectations(t)
	})

	t.Run("Failed Case - Cache Hit (Invalid Data)", func(t *testing.T) {
		_, cache, service := setupTestService(t)
		cacheKey := generateCacheKey(startOfWeek, endOfWeek)
		invalidCachedData := []byte("invalid_data")
		cache.On("Get", cacheKey).Return(invalidCachedData, nil)

		gramTotalCount, err := service.GetGramPlasticStat(startOfWeek, endOfWeek)

		assert.Error(t, err)
		assert.Zero(t, gramTotalCount)

		cache.AssertExpectations(t)
	})

	t.Run("Success Case - Cache Miss", func(t *testing.T) {
		repo, cache, service := setupTestService(t)
		cacheKey := generateCacheKey(startOfWeek, endOfWeek)
		cache.On("Get", cacheKey).Return(nil, errors.New("cache miss"))
		repo.On("GetGramPlasticStat", startOfWeek, endOfWeek).Return(expectedGramTotalCount, nil)
		cache.On("Set", cacheKey, []byte(strconv.FormatUint(expectedGramTotalCount, 10)), 5*time.Second).Return(nil)

		gramTotalCount, err := service.GetGramPlasticStat(startOfWeek, endOfWeek)

		assert.NoError(t, err)
		assert.Equal(t, expectedGramTotalCount, gramTotalCount)

		repo.AssertExpectations(t)
		cache.AssertExpectations(t)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		repo, cache, service := setupTestService(t)
		cacheKey := generateCacheKey(startOfWeek, endOfWeek)
		expectedError := errors.New("gagal menghitung total gram plastik")
		cache.On("Get", cacheKey).Return(nil, errors.New("cache miss"))
		repo.On("GetGramPlasticStat", startOfWeek, endOfWeek).Return(uint64(0), expectedError)

		gramTotalCount, err := service.GetGramPlasticStat(startOfWeek, endOfWeek)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, gramTotalCount)

		repo.AssertExpectations(t)
		cache.AssertExpectations(t)
	})

	t.Run("Failed Case - Cache Set Error", func(t *testing.T) {
		repo, cache, service := setupTestService(t)
		cacheKey := generateCacheKey(startOfWeek, endOfWeek)
		expectedError := errors.New("gagal menyimpan data ke cache")
		cache.On("Get", cacheKey).Return(nil, errors.New("cache miss"))
		repo.On("GetGramPlasticStat", startOfWeek, endOfWeek).Return(expectedGramTotalCount, nil)
		cache.On("Set", cacheKey, []byte(strconv.FormatUint(expectedGramTotalCount, 10)), 5*time.Second).Return(expectedError)

		gramTotalCount, err := service.GetGramPlasticStat(startOfWeek, endOfWeek)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Zero(t, gramTotalCount)

		repo.AssertExpectations(t)
		cache.AssertExpectations(t)
	})
}

func TestGetLatestTransactions(t *testing.T) {
	limit := 5
	mockTransactions := []*entities.OrderModels{
		{
			User:            entities.UserModels{Name: "John"},
			CreatedAt:       time.Now(),
			PaymentStatus:   "Success",
			GrandTotalPrice: uint64(100),
		},
	}

	var responseData []*dto.LastTransactionResponse
	for _, t := range mockTransactions {
		responseData = append(responseData, &dto.LastTransactionResponse{
			Username:      t.User.Name,
			Date:          t.CreatedAt.Format("02-01-2006"),
			PaymentStatus: t.PaymentStatus,
			TotalPrice:    t.GrandTotalPrice,
		})
	}

	t.Run("Success Case", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		repo.On("GetLatestTransactions", limit).Return(mockTransactions, nil)

		response, err := service.GetLatestTransactions(limit)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, len(mockTransactions), len(response))

		assert.Equal(t, len(mockTransactions), len(responseData))
		assert.Equal(t, "John", responseData[0].Username)
		assert.Equal(t, time.Now().Format("02-01-2006"), responseData[0].Date)
		assert.Equal(t, "Success", responseData[0].PaymentStatus)
		assert.Equal(t, uint64(100), responseData[0].TotalPrice)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		repo, _, service := setupTestService(t)
		expectedError := errors.New("gagal mengambil data transaksi terakhir")
		repo.On("GetLatestTransactions", limit).Return(nil, expectedError)

		response, err := service.GetLatestTransactions(limit)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, expectedError.Error())

		repo.AssertExpectations(t)
	})
}
