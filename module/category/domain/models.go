package domain

import "time"

type CategoryModels struct {
	ID           int        `gorm:"column:id;type:int;primaryKey" json:"id"`
	Name         string     `gorm:"column:name;type:varchar(255)" json:"name"`
	TotalProduct int        `gorm:"column:total_product;type:int" json:"total_product"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
}

func (CategoryModels) TableName() string {
	return "category"
}
