package repository

import (
	"errors"
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
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsersPassword(userID uint64) (string, error) {
	var user domain.UserModels
	if err := r.db.Table("users").Select("password").Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Password, nil
}

func (r *UserRepository) ChangePassword(userID uint64, newPasswordHash string) error {
	var user domain.UserModels
	if err := r.db.Model(&user).Where("id = ?", userID).Update("password", newPasswordHash).Error; err != nil {
		return err
	}
	return nil
}
