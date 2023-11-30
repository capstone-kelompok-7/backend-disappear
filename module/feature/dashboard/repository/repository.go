package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/dashboard"
	"gorm.io/gorm"
	"time"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) dashboard.RepositoryDashboardInterface {
	return &DashboardRepository{
		db: db,
	}
}

func (r *DashboardRepository) CountProducts() (int64, error) {
	var count int64
	if err := r.db.Model(&entities.ProductModels{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (r *DashboardRepository) CountUsers() (int64, error) {
	var count int64
	if err := r.db.Model(&entities.UserModels{}).
		Where("is_verified = ? AND role = ?", 1, "customer").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DashboardRepository) CountOrder() (int64, error) {
	var count int64
	if err := r.db.Model(&entities.OrderModels{}).
		Where("order_status = ? AND payment_status = ?", "proses", "konfirmasi").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DashboardRepository) CountIncome() (float64, error) {
	var totalAmount float64

	firstDay := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02")
	lastDay := time.Now().AddDate(0, 1, -time.Now().Day()).Format("2006-01-02")
	if err := r.db.Model(&entities.OrderModels{}).
		Where("order_status = ? AND payment_status = ? AND created_at BETWEEN ? AND ?",
			"proses", "konfirmasi", firstDay, lastDay).
		Select("SUM(total_amount_paid)").
		Scan(&totalAmount).Error; err != nil {
		return 0, err
	}
	return totalAmount, nil
}
