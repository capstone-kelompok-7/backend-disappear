package service

import (
	"errors"
	"math"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/voucher"
)

type VoucherService struct {
	repo        voucher.RepositoryVoucherInterface
	userService users.ServiceUserInterface
}

func NewVoucherService(repo voucher.RepositoryVoucherInterface, userService users.ServiceUserInterface) voucher.ServiceVoucherInterface {
	return &VoucherService{
		repo:        repo,
		userService: userService,
	}
}

func (s *VoucherService) CreateVoucher(newData *entities.VoucherModels) (*entities.VoucherModels, error) {
	if existingVoucher, _ := s.repo.GetVoucherByCode(newData.Code); existingVoucher != nil {
		return nil, errors.New("kode kupon sudah digunakan")
	}

	newVoucher := &entities.VoucherModels{
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

func (s *VoucherService) UpdateVoucher(voucherID uint64, req *entities.VoucherModels) error {
	vouchers, err := s.repo.GetVoucherById(voucherID)
	if err != nil {
		return errors.New("kupon tidak ditemukan")
	}
	if existingVoucher, _ := s.repo.GetVoucherByCode(req.Code); existingVoucher != nil {
		return errors.New("kode kupon sudah digunakan")
	}
	updatedVoucher := &entities.VoucherModels{
		ID:          voucherID,
		Name:        req.Name,
		Code:        req.Code,
		Category:    req.Category,
		Description: req.Description,
		Discount:    req.Discount,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		MinPurchase: req.MinPurchase,
		Stock:       req.Stock,
		Status:      req.Status,
		UpdatedAt:   time.Now(),
	}
	currentTime := time.Now()
	if currentTime.After(updatedVoucher.EndDate) {
		updatedVoucher.Status = "Kadaluwarsa"
	} else {
		updatedVoucher.Status = "Belum Kadaluwarsa"
	}
	err = s.repo.UpdateVoucher(vouchers.ID, updatedVoucher)
	if err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) DeleteVoucher(voucherID uint64) error {
	vouchers, err := s.repo.GetVoucherById(voucherID)
	if err != nil {
		return errors.New("kupon tidak ditemukan")
	}
	if err := s.repo.DeleteVoucher(vouchers.ID); err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) GetVoucherById(voucherID uint64) (*entities.VoucherModels, error) {
	result, err := s.repo.GetVoucherById(voucherID)
	if err != nil {
		return nil, errors.New("kupon tidak ditemukan")
	}

	return result, nil
}

func (s *VoucherService) GetAllVoucher(page, perPage int) ([]*entities.VoucherModels, int64, error) {
	vouchers, err := s.repo.FindAllVoucher(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalVoucherCount()
	if err != nil {
		return nil, 0, err
	}

	return vouchers, totalItems, nil
}

func (s *VoucherService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	pageInt := page
	if pageInt <= 0 {
		pageInt = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	if pageInt > totalPages {
		pageInt = totalPages
	}

	return pageInt, totalPages
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

func (s *VoucherService) CanClaimsVoucher(userID, voucherID uint64) (bool, error) {
	userLevel, err := s.userService.GetUserLevel(userID)
	if err != nil {
		return false, err
	}

	voucherCategory, err := s.repo.GetVoucherCategory(voucherID)
	if err != nil {
		return false, err
	}
	if voucherCategory == "All Customer" {
		return true, nil
	}

	switch userLevel {
	case "Gold":
		return true, nil
	case "Silver":
		if voucherCategory == "Bronze" || voucherCategory == "Silver" {
			return true, nil
		}
	case "Bronze":
		if voucherCategory == "Bronze" {
			return true, nil
		}
	}
	return false, nil
}

func (s *VoucherService) ClaimVoucher(req *entities.VoucherClaimModels) error {
	vouchers, err := s.repo.GetVoucherById(req.VoucherID)
	if err != nil {
		return errors.New("kupon tidak ditemukan")
	}

	if vouchers.Stock == 0 {
		return errors.New("stok kupon sudah habis")
	}

	hasAccess, err := s.CanClaimsVoucher(req.UserID, vouchers.ID)
	if err != nil {
		return err
	}

	if !hasAccess {
		return errors.New("level anda masih belum mencukupi")
	}
	claimed, err := s.repo.IsVoucherAlreadyClaimed(req.UserID, vouchers.ID)
	if err != nil {
		return err
	}

	if claimed {
		return errors.New("kupon telah diklaim")
	}

	newClaims := &entities.VoucherClaimModels{
		UserID:    req.UserID,
		VoucherID: req.VoucherID,
	}
	if err := s.repo.ClaimVoucher(newClaims); err != nil {
		return err
	}
	if err := s.repo.ReduceStockWhenClaimed(vouchers.ID, 1); err != nil {
		return err
	}

	return nil

}

func (s *VoucherService) DeleteVoucherClaims(userID, voucherID uint64) error {
	err := s.repo.DeleteUserVoucherClaims(userID, voucherID)
	if err != nil {
		return err
	}
	return nil

}

func (s *VoucherService) GetUserVouchers(userID uint64) ([]*entities.VoucherClaimModels, error) {
	userVouchers, err := s.repo.GetUserVoucherClaims(userID)
	if err != nil {
		return nil, err
	}
	return userVouchers, nil
}

func (s *VoucherService) GetVoucherByStatus(page, perPage int, status string) ([]*entities.VoucherModels, int64, error) {
	vouchers, err := s.repo.FindByStatus(page, perPage, status)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalVoucherCountByStatus(status)
	if err != nil {
		return nil, 0, err
	}

	return vouchers, totalItems, nil
}

func (s *VoucherService) GetVoucherByCategory(page, perPage int, category string) ([]*entities.VoucherModels, int64, error) {
	vouchers, err := s.repo.FindByCategory(page, perPage, category)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalVoucherCountByCategory(category)
	if err != nil {
		return nil, 0, err
	}

	return vouchers, totalItems, nil
}

func (s *VoucherService) GetVoucherByStatusCategory(page, perPage int, status, category string) ([]*entities.VoucherModels, int64, error) {
	vouchers, err := s.repo.FindByStatusCategory(page, perPage, status, category)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalVoucherCountByStatusCategory(status, category)
	if err != nil {
		return nil, 0, err
	}

	return vouchers, totalItems, nil
}

func (s *VoucherService) GetAllVoucherToClaims(limit int, userID uint64) ([]*entities.VoucherModels, error) {
	vouchers, err := s.repo.FindAllVoucherToClaims(limit, userID)
	if err != nil {
		return nil, err
	}

	filteredVouchers := make([]*entities.VoucherModels, 0)
	for _, voucher := range vouchers {
		claimed, err := s.repo.IsVoucherAlreadyClaimed(userID, voucher.ID)
		if err != nil {
			return nil, err
		}

		if !claimed {
			filteredVouchers = append(filteredVouchers, voucher)
		}
	}

	return filteredVouchers, nil
}
