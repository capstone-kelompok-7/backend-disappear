package service

import (
	"errors"
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
)

type UserService struct {
	repo users.RepositoryUserInterface
}

func NewUserService(repo users.RepositoryUserInterface) users.ServiceUserInterface {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetAllUsers() ([]*domain.UserModels, error) {
	result, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) GetUsersById(userId uint64) (*domain.UserModels, error) {
	result, err := s.repo.GetUsersById(userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) GetUsersByEmail(email string) (*domain.UserModels, error) {
	result, err := s.repo.GetUsersByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}
	if result == nil {
		return nil, errors.New("Pengguna tidak ditemukan")
	}
	return result, nil
}
