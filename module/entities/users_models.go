package entities

import "time"

type UserModels struct {
	ID             uint64          `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	SocialID       string          `gorm:"column:social_id;type:VARCHAR(255)" json:"social_id"`
	Provider       string          `gorm:"column:provider;type:VARCHAR(255)" json:"provider"`
	Email          string          `gorm:"column:email;type:VARCHAR(255)" json:"email"`
	Password       string          `gorm:"column:password;type:VARCHAR(255)" json:"password"`
	Phone          string          `gorm:"column:phone;type:VARCHAR(255)" json:"phone"`
	Role           string          `gorm:"column:role;type:VARCHAR(255)" json:"role"`
	Name           string          `gorm:"column:name;type:VARCHAR(255)" json:"name"`
	PhotoProfile   string          `gorm:"column:photo_profile;type:VARCHAR(255)" json:"photo_profile"`
	TotalGram      uint64          `gorm:"column:total_gram;type:BIGINT UNSIGNED" json:"total_gram"`
	TotalChallenge uint64          `gorm:"column:total_challenge;type:BIGINT UNSIGNED" json:"total_challenge"`
	Level          string          `gorm:"column:level;type:VARCHAR(255)" json:"level"`
	Exp            uint64          `gorm:"column:exp;type:BIGINT UNSIGNED" json:"exp"`
	IsVerified     bool            `gorm:"column:is_verified;default:false" json:"is_verified"`
	LastLogin      time.Time       `gorm:"column:last_login;type:timestamp;default:CURRENT_TIMESTAMP" json:"last_login"`
	DeviceToken    string          `gorm:"column:device_token;type:VARCHAR(255)" json:"device_token"`
	CreatedAt      time.Time       `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      *time.Time      `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
	Address        []AddressModels `gorm:"foreignKey:UserID" json:"addresses"`
	Reviews        []ReviewModels  `gorm:"foreignKey:UserID" json:"reviews"`
}

type AddressModels struct {
	ID           uint64     `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey" json:"id"`
	UserID       uint64     `gorm:"column:user_id;type:BIGINT UNSIGNED" json:"user_id"`
	AcceptedName string     `gorm:"column:accepted_name;type:VARCHAR(255)" json:"accepted_name"`
	Phone        string     `gorm:"column:phone;type:VARCHAR(255)" json:"phone"`
	Address      string     `gorm:"column:address;type:VARCHAR(255)" json:"address"`
	IsPrimary    bool       `gorm:"column:is_primary;type:BOOLEAN" json:"is_primary"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
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
