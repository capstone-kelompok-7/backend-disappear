package entities

import "time"

type VoucherModels struct {
	ID          uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id" `
	Name        string     `gorm:"column:name;type:VARCHAR(255)" json:"name" `
	Code        string     `gorm:"column:code;type:VARCHAR(255)" json:"code" `
	Category    string     `gorm:"column:category;type:VARCHAR(255)" json:"category" `
	Description string     `gorm:"column:description;type:TEXT" json:"description" `
	Discount    uint64     `gorm:"column:discount;type:BIGINT UNSIGNED" json:"discount" `
	StartDate   time.Time  `gorm:"column:start_date;type:DATETIME" json:"start_date" `
	EndDate     time.Time  `gorm:"column:end_date; type:DATETIME" json:"end_date" `
	MinPurchase uint64     `gorm:"column:min_purchase;type:BIGINT UNSIGNED" json:"min_purchase" `
	Stock       uint64     `gorm:"column:stock;type:BIGINT UNSIGNED" json:"stock" `
	Status      string     `gorm:"column:status;type:VARCHAR(255)" json:"status" `
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
}

func (VoucherModels) TableName() string {
	return "vouchers"
}

type VoucherClaimModels struct {
	ID        uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	UserID    uint64         `gorm:"column:user_id;type:BIGINT UNSIGNED" json:"user_id"`
	VoucherID uint64         `gorm:"column:voucher_id;type:BIGINT UNSIGNED" json:"voucher_id"`
	User      *UserModels    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Voucher   *VoucherModels `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
}

func (VoucherClaimModels) TableName() string {
	return "voucher_claims"
}
