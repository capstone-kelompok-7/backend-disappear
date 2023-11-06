package service

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/voucher/domain"
)

type VoucherService struct {
	repo voucher.RepositoryVoucherInterface
}

func NewVoucherService(repo voucher.RepositoryVoucherInterface) voucher.ServiceVoucherInterface {
	return &VoucherService{
		repo: repo,
	}
}

func (s *VoucherService) CreateVoucher(newData domain.VoucherModels) (*domain.VoucherModels, error) {
	result, err := s.repo.CreateVoucher(newData)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *VoucherService) GetAllVouchers(currentPage int, limit int) ([]domain.VoucherModels, error) {
	result, err := s.repo.GetAllVouchers(currentPage, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}
