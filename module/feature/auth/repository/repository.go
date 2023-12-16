package repository

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/auth"

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

func (r *AuthRepository) Register(newData *entities.UserModels) (*entities.UserModels, error) {
	if err := r.db.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (r *AuthRepository) Login(email string) (*entities.UserModels, error) {
	var user entities.UserModels
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) SaveOTP(otp *entities.OTPModels) (*entities.OTPModels, error) {
	err := r.db.Create(&otp).Error
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *AuthRepository) FindValidOTP(userID int, otp string) (*entities.OTPModels, error) {
	var validOTP entities.OTPModels
	err := r.db.Where("user_id = ? AND otp = ? AND expired_otp > ?", userID, otp, time.Now().Unix()).Find(&validOTP).Error
	if err != nil {
		return &validOTP, err
	}

	return &validOTP, nil
}

func (r *AuthRepository) UpdateUser(user *entities.UserModels) (*entities.UserModels, error) {
	err := r.db.Model(&user).Updates(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *AuthRepository) DeleteOTP(otp *entities.OTPModels) error {
	if err := r.db.Delete(&otp).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) DeleteUserOTP(userId uint64) error {
	if err := r.db.Where("user_id = ?", userId).Delete(&entities.OTPModels{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) ResetPassword(email, newPasswordHash string) error {
	var user entities.UserModels
	if err := r.db.Model(&user).Where("email = ?", email).Update("password", newPasswordHash).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) LoginSocial(socialID string) (*entities.UserModels, error) {
	var user entities.UserModels
	if err := r.db.Where("social_id = ?", socialID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserBySocialID(socialID string) (*entities.UserModels, error) {
	var user entities.UserModels
	if err := r.db.Table("users").Where("social_id = ? AND deleted_at IS NULL", socialID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateLastLogin(userID uint64, lastLogin time.Time) error {
	var user entities.UserModels
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	user.LastLogin = lastLogin
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) CekDeviceTokenByEmail(email string) (string, error) {
	var user entities.UserModels
	if err := r.db.Table("users").Where("email = ? AND  deleted_at IS NULL", email).First(&user).Error; err != nil {
		return "", err
	}

	return user.DeviceToken, nil
}

func (r *AuthRepository) UpdateDeviceTokenByID(email string, deviceToken string) (*entities.UserModels, error) {
	var user entities.UserModels
	if err := r.db.Model(&user).Where("email = ?", email).Update("device_token", deviceToken).Error; err != nil {
		return nil, err
	}

	return &user, nil

}
