package repository

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"gorm.io/gorm"
)

type ChallengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) challenge.RepositoryChallengeInterface {
	return &ChallengeRepository{
		db: db,
	}
}

func (r *ChallengeRepository) FindAll(page, perpage int) ([]entities.ChallengeModels, error) {
	var challenge []entities.ChallengeModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Find(&challenge).Error
	if err != nil {
		return challenge, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) GetTotalChallengeCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeModels{}).Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) FindByTitle(page, perpage int, title string) ([]entities.ChallengeModels, error) {
	var challenge []entities.ChallengeModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("title LIKE?", "%"+title+"%").Find(&challenge).Error
	if err != nil {
		return challenge, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) GetTotalChallengeCountByTitle(title string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeModels{}).Where("title LIKE?", "%"+title+"%").Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) CreateChallenge(newData entities.ChallengeModels) (entities.ChallengeModels, error) {
	if err := r.db.Create(&newData).Error; err != nil {
		return newData, err
	}

	return newData, nil
}

func (r *ChallengeRepository) GetChallengeById(id uint64) (entities.ChallengeModels, error) {
	var challenge = entities.ChallengeModels{}
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&challenge).Error; err != nil {
		return challenge, err
	}

	return challenge, nil
}

func (r *ChallengeRepository) UpdateChallenge(challengeID uint64, updatedChallenge entities.ChallengeModels) (entities.ChallengeModels, error) {
	if err := r.db.Model(&entities.ChallengeModels{}).Where("id = ?", challengeID).Updates(updatedChallenge).Error; err != nil {
		return entities.ChallengeModels{}, err
	}

	return updatedChallenge, nil
}

func (r *ChallengeRepository) DeleteChallenge(id uint64) error {
	var challenge entities.ChallengeModels
	if err := r.db.First(&challenge, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	if err := r.db.Model(&challenge).Update("DeletedAt", time.Now()).Error; err != nil {
		return err
	}

	return nil
}
