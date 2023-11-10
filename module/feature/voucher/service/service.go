package service

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
)

type VoucherService struct {
	repo voucher.RepositoryVoucherInterface
}

func NewVoucherService(repo voucher.RepositoryVoucherInterface) voucher.ServiceVoucherInterface {
	return &VoucherService{
		repo: repo,
	}
}

func (s *VoucherService) CreateVoucher(newData entities.VoucherModels) (*entities.VoucherModels, error) {
	result, err := s.repo.CreateVoucher(newData)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *VoucherService) GetAllVouchers(currentPage int, limit int, search string) ([]entities.VoucherModels, error) {
	result, err := s.repo.GetAllVouchers(currentPage, limit, search)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *VoucherService) GetAllVouchersToCalculatePage() ([]entities.VoucherModels, error) {
	result, err := s.repo.GetAllVouchersToCalculatePage()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *VoucherService) EditVoucherById(data entities.VoucherModels) (*entities.VoucherModels, error) {
	result, err := s.repo.EditVoucherById(data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *VoucherService) DeleteVoucherById(id int) error {
	result := s.repo.DeleteVoucherById(id)
	if result != nil {
		return nil
	}
	return nil
}

func (s *VoucherService) GetVoucherById(id int) (*entities.VoucherModels, error) {
	result, err := s.repo.GetVoucherById(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
