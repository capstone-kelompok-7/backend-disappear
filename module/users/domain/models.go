package domain

import "time"

type UserModels struct {
	ID           uint64          `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	Email        string          `gorm:"column:email;type:VARCHAR(255)" json:"email"`
	Password     string          `gorm:"column:password;type:VARCHAR(255)" json:"password"`
	Phone        string          `gorm:"column:phone;type:VARCHAR(255)" json:"phone"`
	Role         string          `gorm:"column:role;type:VARCHAR(255)" json:"role"`
	Name         string          `gorm:"column:name;type:VARCHAR(255)" json:"name"`
	PhotoProfile string          `gorm:"column:photo_profile;type:VARCHAR(255)" json:"photo_profile"`
	TotalGram    float64         `gorm:"column:total_gram;type:DECIMAL(10, 2)" json:"total_gram"`
	IsVerified   bool            `gorm:"column:is_verified;default:false" json:"is_verified"`
	LevelID      int             `gorm:"column:level_id;foreignKey:ID" json:"level_id"`
	Address      []AddressModels `gorm:"foreignKey:UserID" json:"addresses"`
	CreatedAt    time.Time       `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time      `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
}

type AddressModels struct {
	ID           uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	UserID       uint64     `gorm:"column:user_id;type:BIGINT UNSIGNED" json:"user_id"`
	AcceptedName string     `gorm:"column:accepted_name;type:VARCHAR(255)" json:"accepted_name"`
	Street       string     `gorm:"column:street;type:VARCHAR(255)" json:"street"`
	SubDistrict  string     `gorm:"column:sub_district;type:VARCHAR(255)" json:"sub_district"`
	City         string     `gorm:"column:city;type:VARCHAR(255)" json:"city"`
	Province     string     `gorm:"column:province;type:VARCHAR(255)" json:"province"`
	PostalCode   int        `gorm:"column:postal_code;type:INT" json:"postal_code"`
	Note         string     `gorm:"column:note;type:TEXT" json:"note"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:TIMESTAMP;index" json:"deleted_at"`
}

type OTPModels struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id" `
	UserID     int        `gorm:"index;unique" json:"user_id" `
	User       UserModels `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"user" `
	OTP        string     `gorm:"column:otp;type:varchar(255)" json:"otp"`
	ExpiredOTP int64      `gorm:"column:expired_otp;type:bigint" json:"expired_otp" `
}

func (UserModels) TableName() string {
	return "users"
}

func (AddressModels) TableName() string {
	return "address"
}

func (OTPModels) TableName() string {
	return "otp"
}
