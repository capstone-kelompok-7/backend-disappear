package entities

import (
	"time"
)

type ChallengeModels struct {
	ID          uint64                `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	Title       string                `gorm:"column:title;type:varchar(255)" json:"title"`
	Photo       string                `gorm:"column:photo;type:varchar(255)" json:"photo"`
	StartDate   time.Time             `gorm:"column:start_date;type:DATETIME" json:"start_date" `
	EndDate     time.Time             `gorm:"column:end_date;type:DATETIME" json:"end_date" `
	Description string                `gorm:"column:description;type:text" json:"description"`
	Status      string                `gorm:"column:status;type:varchar(255)" json:"status"`
	Exp         uint64                `gorm:"column:exp;type:int" json:"exp"`
	CreatedAt   time.Time             `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time             `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time            `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
	Forms       []ChallengeFormModels `gorm:"foreignKey:ChallengeID" json:"forms"`
}

type ChallengeFormModels struct {
	ID          uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	UserID      uint64     `gorm:"column:user_id;type:BIGINT UNSIGNED" json:"user_id"`
	ChallengeID uint64     `gorm:"column:challenge_id;type:BIGINT UNSIGNED" json:"challenge_id"`
	Username    string     `gorm:"column:username;type:varchar(255)" json:"username"`
	Photo       string     `gorm:"column:photo;type:varchar(255)" json:"photo"`
	Status      string     `gorm:"column:status;type:varchar(255)" json:"status"`
	Exp         uint64     `gorm:"column:exp;type:int" json:"exp"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
}

func (ChallengeModels) TableName() string {
	return "challenges"
}

func (ChallengeFormModels) TableName() string {
	return "challenges_form"
}
