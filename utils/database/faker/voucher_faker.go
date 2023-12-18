package faker

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
)

func VoucherFaker(db *gorm.DB) []*entities.VoucherModels {
	vouchers := []*entities.VoucherModels{
		{
			ID:          1,
			Name:        "Diskon 20% - Belanja Pertama",
			Code:        "FIRST20",
			Category:    "Bronze",
			Description: "Voucher diskon 20% untuk pembelian pertama di toko kami.",
			Discount:    20,
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 1, 0),
			MinPurchase: 50000,
			Stock:       100,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
		{
			ID:          2,
			Name:        "Gratis Ongkir - Belanja di Atas 100Rb",
			Code:        "SHIPFREE",
			Category:    "All Customer",
			Description: "Voucher gratis ongkir untuk pembelian di atas 100.000 rupiah.",
			Discount:    10000,
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 2, 0),
			MinPurchase: 100000,
			Stock:       200,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
	}

	return vouchers
}
