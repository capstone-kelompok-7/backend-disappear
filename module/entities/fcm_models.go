package entities

import "time"

type FcmModels struct {
	ID        uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	UserID    uint64     `gorm:"column:user_id;type:BIGINT UNSIGNED" json:"user_id"`
	Title     string     `gorm:"column:title;type:varchar(255)" json:"title"`
	Body      string     `gorm:"column:body;type:text" json:"body"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
}

func (FcmModels) TableName() string {
	return "fcms"
}
