package repository

import (
	"errors"
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) users.RepositoryUserInterface {
	return &UserRepository{
		db: db,
	}
}
func (r *UserRepository) GetUsersById(userId uint64) (*domain.UserModels, error) {
	var user domain.UserModels
	if err := r.db.Where("id", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, errors.New("id not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]*domain.UserModels, error) {
	var users []*domain.UserModels
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUsersByEmail(email string) (*domain.UserModels, error) {
	var user domain.UserModels
	if err := r.db.Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) ChangePassword(password string) (*domain.UserModels, error) {
	var user domain.UserModels
	if err := r.db.Table("users").Where("password = ?", password).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdatePassword(userID uint64, newPassword string) error {
	return r.db.Model(&domain.UserModels{}).
		Where("id = ?", userID).
		Update("password", newPassword).
		Error
}

func (r *UserRepository) ComparePassword(oldPass string) (*domain.UserModels, error) {
	var user domain.UserModels
	if err := r.db.Where("password = ?", oldPass).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
