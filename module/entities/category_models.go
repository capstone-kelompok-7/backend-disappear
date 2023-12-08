package entities

import "time"

type CategoryModels struct {
	ID           uint64          `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	Name         string          `gorm:"column:name;type:varchar(255)" json:"name"`
	Photo        string          `gorm:"column:photo;type:varchar(255)" json:"photo"`
	TotalProduct uint64          `gorm:"column:total_product;type:BIGINT UNSIGNED" json:"total_product"`
	CreatedAt    time.Time       `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time      `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
	Products     []ProductModels `gorm:"many2many:product_categories;" json:"products"`
}

func (CategoryModels) TableName() string {
	return "category"
}
