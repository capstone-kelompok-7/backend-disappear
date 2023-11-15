package entities

import "time"

type ReviewModels struct {
	ID          int        `gorm:"column:id;type:int;primaryKey" json:"id"`
	ProductID   int        `gorm:"column:product_id;type:int" json:"product_id"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	Date        time.Time  `gorm:"column:date;type:date" json:"date"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
}

func (ReviewModels) TableName() string {
	return "reviews"
}
