package repository

import (
	"errors"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users/dto"
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

func (r *UserRepository) GetUsersById(userId uint64) (*entities.UserModels, error) {
	var user entities.UserModels
	if err := r.db.Preload("Address").Where("id = ? AND deleted_at IS NULL", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsersByEmail(email string) (*entities.UserModels, error) {
	var user entities.UserModels
	if err := r.db.Table("users").Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsersPassword(userID uint64) (string, error) {
	var user entities.UserModels
	if err := r.db.Table("users").Select("password").Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Password, nil
}

func (r *UserRepository) ChangePassword(userID uint64, newPasswordHash string) error {
	var user entities.UserModels
	if err := r.db.Model(&user).Where("id = ?", userID).Update("password", newPasswordHash).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindAll(page, perPage int) ([]*entities.UserModels, error) {
	var user []*entities.UserModels
	offset := (page - 1) * perPage
	err := r.db.Offset(offset).Limit(perPage).Preload("Address").Where("deleted_at IS NULL").Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByName(page, perPage int, name string) ([]*entities.UserModels, error) {
	var user []*entities.UserModels
	offset := (page - 1) * perPage
	query := r.db.Offset(offset).Limit(perPage).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetTotalUserCountByName(name string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.UserModels{}).Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *UserRepository) GetTotalUserCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.UserModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *UserRepository) EditProfile(userID uint64, updatedData dto.EditProfileRequest) (*entities.UserModels, error) {
	var user *entities.UserModels
	if err := r.db.Model(&user).Where("id = ?", userID).Updates(&updatedData).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) DeleteAccount(userID uint64) error {
	user := &entities.UserModels{}
	if err := r.db.First(user, userID).Error; err != nil {
		return err
	}

	if err := r.db.Model(user).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUserExp(userID uint64, exp uint64) (*entities.UserModels, error) {
	user := &entities.UserModels{}
	if err := r.db.Model(user).Where("id = ?", userID).Update("exp", exp).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUserChallengeFollow(userID uint64, totalChallenge uint64) (*entities.UserModels, error) {
	user := &entities.UserModels{}
	if err := r.db.Model(user).Where("id = ?", userID).Update("total_challenge", totalChallenge).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUserContribution(userID uint64, gramPlastic uint64) (*entities.UserModels, error) {
	user := &entities.UserModels{}
	if err := r.db.Model(user).Where("id = ?", userID).Update("total_gram", gramPlastic).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUserLevel(userID uint64, level string) error {
	var user entities.UserModels
	if err := r.db.Model(&user).Where("id = ?", userID).Update("level", level).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserLevel(userID uint64) (string, error) {
	var user entities.UserModels
	if err := r.db.Select("level").Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Level, nil
}

func (r *UserRepository) GetFilterLevel(level string) ([]*entities.UserModels, error) {
	var user []*entities.UserModels
	if err := r.db.Where("level = ? AND deleted_at IS NULL", level).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetLeaderboardByExp(limit int) ([]*entities.UserModels, error) {
	var user []*entities.UserModels
	if err := r.db.Where("role = ?", "customer").Order("exp DESC").Limit(limit).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserTransactionActivity(userID uint64) (int, int, int, error) {
	var success []*entities.OrderModels
	var failed []*entities.OrderModels

	if err := r.db.Where("user_id = ?", userID).
		Where("payment_status = ?", "Konfirmasi").
		Where("deleted_at IS NULL").
		Order("created_at desc").
		Find(&success).Error; err != nil {
		return 0, 0, 0, err
	}

	if err := r.db.Where("user_id = ?", userID).
		Where("payment_status = ?", "Gagal").
		Where("deleted_at IS NULL").
		Order("created_at desc").
		Find(&failed).Error; err != nil {
		return 0, 0, 0, err
	}

	numSuccess := len(success)
	numFailed := len(failed)
	total := numSuccess + numFailed

	return numSuccess, numFailed, total, nil
}

func (r *UserRepository) GetUserChallengeActivity(userID uint64) (int, int, int, error) {
	var success []*entities.ChallengeFormModels
	var failed []*entities.ChallengeFormModels

	if err := r.db.Where("user_id = ?", userID).
		Where("status = ?", "Valid").
		Where("deleted_at IS NULL").
		Order("created_at desc").
		Find(&success).Error; err != nil {
		return 0, 0, 0, err
	}

	if err := r.db.Where("user_id = ?", userID).
		Where("status = ?", "Tidak Valid").
		Where("deleted_at IS NULL").
		Order("created_at desc").
		Find(&failed).Error; err != nil {
		return 0, 0, 0, err
	}

	numSuccess := len(success)
	numFailed := len(failed)
	total := numSuccess + numFailed

	return numSuccess, numFailed, total, nil

}

func (r *UserRepository) FindAllPersonalization() ([]*entities.PersonalizationModels, error) {
	var personalizations []*entities.PersonalizationModels

	err := r.db.Preload("Isu").Preload("Category").Preload("User").Find(&personalizations).Error
	if err != nil {
		return nil, err
	}

	return personalizations, nil
}

func (r *UserRepository) FindAllEnvironmentsIsues() ([]*entities.EnvironmentIssuesModels, error) {
	var environmentIssues []*entities.EnvironmentIssuesModels
	err := r.db.Find(&environmentIssues).Error
	if err != nil {
		return nil, err
	}
	return environmentIssues, nil
}

func (r *UserRepository) CreateUserPersonalization(userID uint64, req *dto.UserPersonalizationRequest) ([]*entities.PersonalizationModels, error) {
	var userPersonalizedData []*entities.PersonalizationModels

	for _, isuID := range req.IsuID {
		for _, categoryID := range req.CategoryID {
			newUserPersonalized := &entities.PersonalizationModels{
				UserID:     userID,
				IsuID:      isuID,
				CategoryID: categoryID,
			}

			if err := r.db.Create(newUserPersonalized).Error; err != nil {
				return nil, err
			}

			userPersonalizedData = append(userPersonalizedData, newUserPersonalized)
		}
	}

	return userPersonalizedData, nil
}

func (r *UserRepository) GetUserPersonalization(userID uint64) ([]*entities.PersonalizationModels, error) {
	var personalizations []*entities.PersonalizationModels
	if err := r.db.Where("user_id = ?", userID).Find(&personalizations).Error; err != nil {
		return nil, err
	}
	return personalizations, nil
}
