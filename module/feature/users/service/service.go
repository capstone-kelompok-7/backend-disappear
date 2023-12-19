package service

import (
	"errors"
	"math"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users/dto"
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

func (s *UserService) GetUsersById(userId uint64) (*entities.UserModels, error) {
	result, err := s.repo.GetUsersById(userId)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	return result, nil
}

func (s *UserService) GetUsersByEmail(email string) (*entities.UserModels, error) {
	result, err := s.repo.GetUsersByEmail(email)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	return result, nil
}

func (s *UserService) ValidatePassword(userID uint64, oldPassword, newPassword, confirmPassword string) error {
	if oldPassword == newPassword {
		return errors.New("Password baru tidak boleh sama dengan password lama")
	}

	if newPassword != confirmPassword {
		return errors.New("Password baru dan konfirmasi password tidak cocok")
	}

	storedPassword, err := s.repo.GetUsersPassword(userID)
	if err != nil {
		return errors.New("Password lama tidak valid")
	}

	isValidOldPassword, err := s.hash.ComparePassword(storedPassword, oldPassword)
	if err != nil || !isValidOldPassword {
		return errors.New("Password lama tidak valid")
	}

	return nil
}

func (s *UserService) ChangePassword(userID uint64, updateRequest dto.UpdatePasswordRequest) error {
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

func (s *UserService) GetAllUsers(page, perPage int) ([]*entities.UserModels, int64, error) {
	user, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalUserCount()
	if err != nil {
		return nil, 0, err
	}

	return user, totalItems, nil
}

func (s *UserService) GetUsersByName(page int, perPage int, name string) ([]*entities.UserModels, int64, error) {
	user, err := s.repo.FindByName(page, perPage, name)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalUserCountByName(name)
	if err != nil {
		return nil, 0, err
	}

	return user, totalItems, nil
}

func (s *UserService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if page > totalPages {
		page = totalPages
	}

	return page, totalPages
}

func (s *UserService) GetNextPage(currentPage int, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}

	return totalPages
}

func (s *UserService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}

	return 1
}

func (s *UserService) EditProfile(userID uint64, updatedData dto.EditProfileRequest) (*entities.UserModels, error) {
	_, err := s.repo.GetUsersById(userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	result, err := s.repo.EditProfile(userID, updatedData)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) DeleteAccount(userID uint64) error {
	_, err := s.repo.GetUsersById(userID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}
	if err := s.repo.DeleteAccount(userID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUserExp(userID uint64, exp uint64) (*entities.UserModels, error) {
	user, err := s.repo.GetUsersById(userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	user.Exp = exp
	updatedUser, err := s.repo.UpdateUserExp(userID, exp)
	if err != nil {
		return nil, err
	}

	level := determineLevel(updatedUser.Exp)
	if level != updatedUser.Level {
		updatedUser.Level = level
		if err := s.repo.UpdateUserLevel(userID, level); err != nil {
			return nil, err
		}
	}
	return updatedUser, nil
}

func (s *UserService) UpdateUserContribution(userID uint64, gramPlastic uint64) (*entities.UserModels, error) {
	user, err := s.repo.GetUsersById(userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	user.TotalGram = gramPlastic
	updatedUser, err := s.repo.UpdateUserContribution(userID, gramPlastic)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func determineLevel(exp uint64) string {
	if exp <= 500 {
		return "Bronze"
	} else if exp <= 1000 {
		return "Silver"
	} else {
		return "Gold"
	}
}

func (s *UserService) GetUserLevel(userID uint64) (string, error) {
	user, err := s.repo.GetUsersById(userID)
	if err != nil {
		return "", errors.New("pengguna tidak ditemukan")
	}
	userLevel, err := s.repo.GetUserLevel(user.ID)
	if err != nil {
		return "", err
	}
	return userLevel, nil
}

func (s *UserService) UpdateUserChallengeFollow(userID uint64, totalChallenge uint64) (*entities.UserModels, error) {
	user, err := s.repo.UpdateUserChallengeFollow(userID, totalChallenge)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetLeaderboardByExp(limit int) ([]*entities.UserModels, error) {
	user, err := s.repo.GetLeaderboardByExp(limit)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserTransactionActivity(userID uint64) (int, int, int, error) {
	user, err := s.repo.GetUsersById(userID)
	if err != nil {
		return 0, 0, 0, errors.New("pengguna tidak ditemukan")
	}

	successOrder, failedOrder, totalOrder, err := s.repo.GetUserTransactionActivity(user.ID)
	if err != nil {
		return 0, 0, 0, err
	}

	return successOrder, failedOrder, totalOrder, nil
}

func (s *UserService) GetUserChallengeActivity(userID uint64) (int, int, int, error) {
	user, err := s.repo.GetUsersById(userID)
	if err != nil {
		return 0, 0, 0, errors.New("pengguna tidak ditemukan")
	}
	successChallenge, failedChallenge, totalChallenge, err := s.repo.GetUserChallengeActivity(user.ID)
	if err != nil {
		return 0, 0, 0, err
	}

	return successChallenge, failedChallenge, totalChallenge, nil
}

func (s *UserService) GetUserProfile(userID uint64) (*entities.UserModels, error) {
	result, err := s.repo.GetUsersById(userID)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (s *UserService) GetUsersBySearchAndFilter(page, perPage int, search, levelFilter string) ([]*entities.UserModels, int64, error) {
	return s.repo.GetAllUsersBySearchAndFilter(page, perPage, search, levelFilter)
}

func (s *UserService) GetUsersByLevel(page, perPage int, level string) ([]*entities.UserModels, int64, error) {
	return s.repo.GetFilterLevel(page, perPage, level)
}
