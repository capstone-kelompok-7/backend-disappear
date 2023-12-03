package service

import (
	"errors"
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils/caching"
	"strconv"
	"time"
)

type DashboardService struct {
	repo  dashboard.RepositoryDashboardInterface
	cache caching.CacheRepository
}

func NewDashboardService(repo dashboard.RepositoryDashboardInterface, cache caching.CacheRepository) dashboard.ServiceDashboardInterface {
	return &DashboardService{
		repo:  repo,
		cache: cache,
	}
}
func generateCacheKey(startOfWeek, endOfWeek time.Time) string {
	return fmt.Sprintf("gramPlasticStat:%s_%s", startOfWeek.Format(time.RFC3339), endOfWeek.Format(time.RFC3339))
}

func (s *DashboardService) GetCardDashboard() (int64, int64, int64, float64, error) {
	productCount, err := s.repo.CountProducts()
	if err != nil {
		return 0, 0, 0, 0.0, errors.New("gagal menghitung total produk")
	}

	orderCount, err := s.repo.CountOrder()
	if err != nil {
		return 0, 0, 0, 0.0, errors.New("gagal menghitung total pesanan")
	}

	userCount, err := s.repo.CountUsers()
	if err != nil {
		return 0, 0, 0, 0.0, errors.New("gagal menghitung total pelanggan")
	}
	inComeCount, err := s.repo.CountIncome()
	if err != nil {
		return 0, 0, 0, 0.0, errors.New("gagal menghitung total pendapatan")
	}

	return productCount, userCount, orderCount, inComeCount, nil
}

func (s *DashboardService) GetLandingPage() (int64, int64, int64, error) {
	userCount, err := s.repo.CountUsers()
	if err != nil {
		return 0, 0, 0, err
	}

	gramPlastic, err := s.repo.CountTotalGram()
	if err != nil {
		return 0, 0, 0, err
	}

	orderCount, err := s.repo.CountOrder()
	if err != nil {
		return 0, 0, 0, err
	}

	return userCount, gramPlastic, orderCount, nil
}

func (s *DashboardService) GetProductReviewsWithMaxTotal() ([]*entities.ProductModels, error) {
	productReviews, err := s.repo.GetProductWithMaxReviews()
	if err != nil {
		return nil, err
	}
	return productReviews, nil
}

func (s *DashboardService) GetGramPlasticStat(startOfWeek, endOfWeek time.Time) (uint64, error) {
	cacheKey := generateCacheKey(startOfWeek, endOfWeek)
	cachedData, err := s.cache.Get(cacheKey)
	if err == nil {
		gramTotalCount, parseErr := strconv.ParseUint(string(cachedData), 10, 64)
		if parseErr != nil {
			return 0, parseErr
		}
		return gramTotalCount, nil
	}
	gramTotalCount, err := s.repo.GetGramPlasticStat(startOfWeek, endOfWeek)
	if err != nil {
		return 0, errors.New("gagal menghitung total gram plastik")
	}
	err = s.cache.Set(cacheKey, []byte(strconv.FormatUint(gramTotalCount, 10)), 5*time.Second)
	if err != nil {
		return 0, err
	}
	return gramTotalCount, nil
}

func (s *DashboardService) GetLatestTransactions(limit int) ([]*dto.LastTransactionResponse, error) {
	transactions, err := s.repo.GetLatestTransactions(limit)
	if err != nil {
		return nil, errors.New("gagal mengambil data transaksi terakhir")
	}
	var responseData []*dto.LastTransactionResponse
	for _, t := range transactions {
		responseData = append(responseData, &dto.LastTransactionResponse{
			Username:      t.User.Name,
			Date:          t.CreatedAt.Format("02-01-2006"),
			PaymentStatus: t.PaymentStatus,
			TotalPrice:    t.GrandTotalPrice,
		})
	}
	return responseData, nil
}
