package entities

import "time"

type ReviewModels struct {
	ID          uint64              `gorm:"column:id;type:int;primaryKey" json:"id"`
	UserID      uint64              `gorm:"column:user_id;type:BIGINT UNSIGNED" json:"user_id"`
	User        UserModels          `gorm:"foreignKey:UserID" json:"user"`
	ProductID   uint64              `gorm:"column:product_id;type:BIGINT UNSIGNED" json:"product_id"`
	Rating      uint64              `gorm:"column:rating;type:BIGINT UNSIGNED" json:"rating"`
	Description string              `gorm:"column:description;type:text" json:"description"`
	Date        time.Time           `gorm:"column:date;type:DATETIME" json:"date"`
	CreatedAt   time.Time           `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time           `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time          `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
	Photos      []ReviewPhotoModels `gorm:"foreignKey:ReviewID" json:"photos"`
}

type ReviewPhotoModels struct {
	ID        uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	ReviewID  uint64     `gorm:"column:review_id;type:BIGINT UNSIGNED" json:"review_id"`
	ImageURL  string     `gorm:"column:url;type:varchar(255)" json:"url"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
}

type ReviewDetail struct {
	Name         string    `json:"name"`
	PhotoProfile string    `json:"photo_profile"`
	Rating       uint64    `json:"rating"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
}

func (ReviewModels) TableName() string {
	return "reviews"
}

func (ReviewPhotoModels) TableName() string {
	return "review_photos"
}
