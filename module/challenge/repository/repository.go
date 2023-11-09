package repository

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/challenge/domain"
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

func (r *ChallengeRepository) FindAll(page, perpage int) ([]domain.ChallengeModels, error) {
	var challenge []domain.ChallengeModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Find(&challenge).Error
	if err != nil {
		return challenge, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) GetTotalChallengeCount() (int64, error) {
	var count int64
	err := r.db.Model(&domain.ChallengeModels{}).Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) FindByTitle(page, perpage int, title string) ([]domain.ChallengeModels, error) {
	var challenge []domain.ChallengeModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("title LIKE?", "%"+title+"%").Find(&challenge).Error
	if err != nil {
		return challenge, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) GetTotalChallengeCountByTitle(title string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.ChallengeModels{}).Where("title LIKE?", "%"+title+"%").Count(&count).Error
	return count, err
}
