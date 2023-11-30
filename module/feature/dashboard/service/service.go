package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard"
)

type DashboardService struct {
	repo dashboard.RepositoryDashboardInterface
}

func NewDashboardService(repo dashboard.RepositoryDashboardInterface) dashboard.ServiceDashboardInterface {
	return &DashboardService{
		repo: repo,
	}
}

func (s *DashboardService) GetCardDashboard() (int64, int64, int64, float64, error) {
	productCount, err := s.repo.CountProducts()
	if err != nil {
		return 0, 0, 0, 0.0, errors.New("Gagal menghitung total produk")
	}

	orderCount, err := s.repo.CountOrder()
	if err != nil {
		return 0, 0, 0, 0.0, errors.New("Gagal menghitung total pesanan")
	}

	userCount, err := s.repo.CountUsers()
	if err != nil {
		return 0, 0, 0, 0.0, errors.New("Gagal menghitung total pelanggan")
	}
	inComeCount, err := s.repo.CountIncome()
	if err != nil {
		return 0, 0, 0, 0.0, errors.New("Gagal menghitung total pendapatan")
	}

	return productCount, userCount, orderCount, inComeCount, nil

}
