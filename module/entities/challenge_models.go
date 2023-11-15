package entities

import (
	"time"
)

type ChallengeModels struct {
	ID          uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	Title       string     `gorm:"column:title;type:varchar(255)" json:"title"`
	Photo       string     `gorm:"column:photo;type:varchar(255)" json:"photo"`
	StartDate   time.Time  `gorm:"column:start_date;type:DATE" json:"start_date"`
	EndDate     time.Time  `gorm:"column:end_date;type:DATE" json:"end_date"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	Winner      string     `gorm:"column:winner;type:varchar(255)" json:"winner"`
	Status      string     `gorm:"column:status;type:varchar(255)" json:"status"`
	Exp         uint64     `gorm:"column:exp;type:int" json:"exp"`
	// CreatedAt   time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	// UpdatedAt   time.Time  `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	// DeletedAt   *time.Time `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
	CreatedAt time.Time  `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (ChallengeModels) TableName() string {
	return "challenges"
}
