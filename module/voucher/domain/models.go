package domain

import "time"

type VoucherModels struct {
	ID          uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	Name        string     `gorm:"column:name;type:VARCHAR(255)" json:"name"`
	Code        string     `gorm:"column:code;type:VARCHAR(255)" json:"code"`
	Category    string     `gorm:"column:category;type:VARCHAR(255)" json:"category"`
	Description string     `gorm:"column:description;type:TEXT" json:"description"`
	Discouunt   int        `gorm:"column:discount;type:INT" json:"discount"`
	StartDate   string     `gorm:"column:start_date;type:DATETIME" json:"start_date"`
	EndDate     string     `gorm:"column:end_date; type:DATETIME" json:"end_date"`
	MinAmount   float64    `gorm:"column:min_amount;type:DECIMAL(10, 2)" json:"min_amount"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:DATETIME;default:NULL;index" json:"deleted_at"`
}

func (VoucherModels) TableName() string {
	return "vouchers"
}
