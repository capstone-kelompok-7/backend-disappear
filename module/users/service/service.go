package service

import (
	"errors"
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

func (s *UserService) ValidatePassword(userID uint64, oldPassword, newPassword, confirmPassword string) error {
	storedPassword, err := s.repo.GetUsersPassword(userID)
	if err != nil {
		return errors.New("Password lama tidak valid")
	}

	isValidOldPassword, err := s.hash.ComparePassword(storedPassword, oldPassword)
	if err != nil || !isValidOldPassword {
		return errors.New("Password lama tidak valid")
	}

	if oldPassword == newPassword {
		return errors.New("Password baru tidak boleh sama dengan password lama")
	}

	if newPassword != confirmPassword {
		return errors.New("Password baru dan konfirmasi password tidak cocok")
	}

	return nil
}

func (s *UserService) ChangePassword(userID uint64, updateRequest domain.UpdatePasswordRequest) error {
	user, err := s.repo.GetUsersById(userID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}
	newPasswordHash, err := s.hash.GenerateHash(updateRequest.NewPassword)
	if err != nil {
		return errors.New("gagal hash password")
	}
	err = s.repo.ChangePassword(user.ID, newPasswordHash)
	if err != nil {
		return err
	}

	return nil
}
