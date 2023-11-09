package repository

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/auth"
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) auth.RepositoryAuthInterface {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Register(newData *domain.UserModels) (*domain.UserModels, error) {
	if err := r.db.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (r *AuthRepository) Login(email string) (*domain.UserModels, error) {
	var user domain.UserModels
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) SaveOTP(otp *domain.OTPModels) (*domain.OTPModels, error) {
	err := r.db.Create(&otp).Error
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *AuthRepository) FindValidOTP(userID int, otp string) (*domain.OTPModels, error) {
	var validOTP domain.OTPModels
	err := r.db.Where("user_id = ? AND otp = ? AND expired_otp > ?", userID, otp, time.Now().Unix()).Find(&validOTP).Error
	if err != nil {
		return &validOTP, err
	}

	return &validOTP, nil
}

func (r *AuthRepository) UpdateUser(user *domain.UserModels) (*domain.UserModels, error) {
	err := r.db.Model(&user).Updates(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *AuthRepository) DeleteOTP(otp *domain.OTPModels) error {
	if err := r.db.Delete(&otp).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) DeleteUserOTP(userId uint64) error {
	if err := r.db.Where("user_id = ?", userId).Delete(&domain.OTPModels{}).Error; err != nil {
		return err
	}

	return nil
}
