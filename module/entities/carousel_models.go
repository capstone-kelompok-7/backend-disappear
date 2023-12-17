package entities

import "time"

type CarouselModels struct {
	ID        uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	Name      string     `gorm:"column:name;type:VARCHAR(255)" json:"name"`
	Photo     string     `gorm:"column:photo;type:VARCHAR(255)" json:"photo"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
}

func (CarouselModels) TableName() string {
	return "carousel"
}
