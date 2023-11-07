package service

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/users/domain"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
)

type UserService struct {
	repo users.RepositoryUserInterface
	hash utils.HashInterface
}

func NewUserService(repo users.RepositoryUserInterface, hash utils.HashInterface) users.ServiceUserInterface {
	return &UserService{
		repo: repo,
		hash: hash,
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
		return nil, err
	}
	return result, nil
}

func (s *UserService) ChangePassword(email, oldPass, newPass string) (*domain.UserModels, error) {
	user, err := s.repo.GetUsersByEmail(email)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.ComparePassword(oldPass)
	if err != nil {
		return nil, err
	}

	newPasswordHash, err := s.hash.GenerateHash(newPass)
	if err != nil {
		return nil, err
	}

	user.Password = newPasswordHash
	result, err := s.repo.ChangePassword(user.Password)
	if err != nil {
		return nil, err
	}
	return result, nil
}
