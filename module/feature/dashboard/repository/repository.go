package repository

import (
	"database/sql"
	"errors"
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
	var totalAmount sql.NullFloat64

	firstDay := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02")
	lastDay := time.Now().AddDate(0, 1, -time.Now().Day()).Format("2006-01-02")
	if err := r.db.Model(&entities.OrderModels{}).
		Where("order_status = ? AND payment_status = ? AND created_at BETWEEN ? AND ?",
			"proses", "konfirmasi", firstDay, lastDay).
		Select("SUM(total_amount_paid)").
		Row().Scan(&totalAmount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0.0, nil
		}
		return 0.0, err
	}
	return totalAmount.Float64, nil
}

func (r *DashboardRepository) CountTotalGram() (int64, error) {
	var totalGramPlastic int64
	if err := r.db.Model(&entities.OrderModels{}).
		Where("payment_status = ? AND order_status = ?", "konfirmasi", "proses").
		Select("COALESCE(SUM(grand_total_gram_plastic), 0)").
		Scan(&totalGramPlastic).Error; err != nil {
		return 0, err
	}
	return totalGramPlastic, nil
}

func (r *DashboardRepository) GetProductWithMaxReviews() ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels
	if err := r.db.
		Preload("ProductReview", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc").Limit(5)
		}).
		Preload("ProductReview.User").
		Where("deleted_at IS NULL").
		Order("total_review desc").
		Limit(1).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *DashboardRepository) GetGramPlasticStat(startOfWeek, endOfWeek time.Time) (uint64, error) {
	var gramTotalCount uint64
	if err := r.db.Model(&entities.OrderModels{}).
		Where("payment_status = 'konfirmasi'").
		Where("created_at BETWEEN ? AND ?", startOfWeek, endOfWeek).
		Select("COALESCE(SUM(grand_total_gram_plastic), 0)").
		Scan(&gramTotalCount).
		Error; err != nil {
		return 0, err
	}
	return gramTotalCount, nil
}

func (r *DashboardRepository) GetLatestTransactions(limit int) ([]*entities.OrderModels, error) {
	var transactions []*entities.OrderModels
	if err := r.db.Order("created_at DESC").
		Limit(limit).
		Preload("User").
		Preload("Address").
		Preload("Voucher").
		Preload("OrderDetails").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
