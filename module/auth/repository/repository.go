package repository

import (
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
