package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher/dto"
	"math"
	"time"
)

type VoucherService struct {
	repo voucher.RepositoryVoucherInterface
}

func NewVoucherService(repo voucher.RepositoryVoucherInterface) voucher.ServiceVoucherInterface {
	return &VoucherService{
		repo: repo,
	}
}

func (s *VoucherService) CreateVoucher(newData entities.VoucherModels) (entities.VoucherModels, error) {
	newVoucher := entities.VoucherModels{
		Name:        newData.Name,
		Code:        newData.Code,
		Category:    newData.Category,
		Description: newData.Description,
		Discount:    newData.Discount,
		StartDate:   newData.StartDate,
		EndDate:     newData.EndDate,
		MinPurchase: newData.MinPurchase,
		Stock:       newData.Stock,
	}
	currentTime := time.Now()
	if currentTime.After(newVoucher.EndDate) {
		newVoucher.Status = "Kadaluwarsa"
	} else {
		newVoucher.Status = "Belum Kadaluwarsa"
	}
	result, err := s.repo.CreateVoucher(newVoucher)
	if err != nil {
		return result, errors.New("gagal menambahkan kupon")
	}
	return result, nil
}

func (s *VoucherService) UpdateVoucher(id uint64, updatedData dto.UpdateVoucherRequest) (entities.VoucherModels, error) {
	_, err := s.repo.GetVoucherById(id)
	if err != nil {
		return entities.VoucherModels{}, errors.New("kupon tidak ditemukan")
	}
	result, err := s.repo.UpdateVoucher(id, updatedData)
	if err != nil {
		return result, errors.New("gagal memperbarui kupon")
	}
	return result, nil
}

func (s *VoucherService) DeleteVoucher(id uint64) error {
	_, err := s.repo.GetVoucherById(id)
	if err != nil {
		return errors.New("kupon tidak ditemukan")
	}
	if err := s.repo.DeleteVoucher(id); err != nil {
		return errors.New("gagal menghapus kupon")
	}
	return nil
}

func (s *VoucherService) GetVoucherById(id uint64) (entities.VoucherModels, error) {
	result, err := s.repo.GetVoucherById(id)
	if err != nil {
		return result, errors.New("kupon tidak ditemukan")
	}

	return result, nil
}

func (s *VoucherService) GetAllVoucher(page, perPage int) ([]entities.VoucherModels, int64, error) {
	vouchers, err := s.repo.FindAllVoucher(page, perPage)
	if err != nil {
		return vouchers, 0, err
	}

	totalItems, err := s.repo.GetTotalVoucherCount()
	if err != nil {
		return vouchers, 0, err
	}

	return vouchers, totalItems, nil
}

func (s *VoucherService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	pageInt := page
	if pageInt <= 0 {
		pageInt = 1
	}

	total_pages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	if pageInt > total_pages {
		pageInt = total_pages
	}

	return pageInt, total_pages
}

func (s *VoucherService) GetNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}
	return totalPages
}

func (s *VoucherService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}
	return 1
}

func (s *VoucherService) GetVouchersByName(page, perPage int, name string) ([]entities.VoucherModels, int64, error) {
	vouchers, err := s.repo.FindVoucherByName(page, perPage, name)
	if err != nil {
		return vouchers, 0, err
	}

	totalItems, err := s.repo.GetTotalVoucherCountByName(name)
	if err != nil {
		return vouchers, 0, err
	}

	return vouchers, totalItems, nil
}
