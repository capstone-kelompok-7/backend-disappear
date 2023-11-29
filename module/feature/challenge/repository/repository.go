package repository

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/dto"
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

func (r *ChallengeRepository) FindAll(page, perpage int) ([]*entities.ChallengeModels, error) {
	var challenge []*entities.ChallengeModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("deleted_at IS NULL").Find(&challenge).Error
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

func (r *ChallengeRepository) FindByTitle(page, perpage int, title string) ([]*entities.ChallengeModels, error) {
	var challenge []*entities.ChallengeModels
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

func (r *ChallengeRepository) FindByStatus(page, perpage int, status string) ([]*entities.ChallengeModels, error) {
	var challenge []*entities.ChallengeModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("status = ?", status).Find(&challenge).Error
	if err != nil {
		return challenge, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) GetTotalChallengeCountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeModels{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) CreateChallenge(newData *entities.ChallengeModels) (*entities.ChallengeModels, error) {
	if err := r.db.Create(&newData).Error; err != nil {
		return newData, err
	}

	return newData, nil
}

func (r *ChallengeRepository) GetChallengeById(id uint64) (*entities.ChallengeModels, error) {
	var challenge = &entities.ChallengeModels{}
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&challenge).Error; err != nil {
		return challenge, err
	}

	return challenge, nil
}

func (r *ChallengeRepository) UpdateChallenge(challengeID uint64, updatedChallenge *entities.ChallengeModels) (*entities.ChallengeModels, error) {
	if err := r.db.Model(&entities.ChallengeModels{}).Where("id = ?", challengeID).Updates(updatedChallenge).Error; err != nil {
		return &entities.ChallengeModels{}, err
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

func (r *ChallengeRepository) CreateSubmitChallengeForm(form *entities.ChallengeFormModels) (*entities.ChallengeFormModels, error) {
	err := r.db.Create(&form).Error
	if err != nil {
		return nil, err
	}
	return form, nil
}

func (r *ChallengeRepository) GetAllSubmitChallengeForm(page, perpage int) ([]*entities.ChallengeFormModels, error) {
	var participants []*entities.ChallengeFormModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("deleted_at IS NULL").Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *ChallengeRepository) GetSubmitChallengeFormByStatus(page, perpage int, status string) ([]*entities.ChallengeFormModels, error) {
	var participants []*entities.ChallengeFormModels
	offset := (page - 1) * perpage
	err := r.db.Where("status = ?", status).Offset(offset).Limit(perpage).Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *ChallengeRepository) GetTotalSubmitChallengeFormCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeFormModels{}).Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) GetTotalSubmitChallengeFormCountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeFormModels{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) GetSubmitChallengeFormById(id uint64) (*entities.ChallengeFormModels, error) {
	var participant = &entities.ChallengeFormModels{}
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&participant).Error; err != nil {
		return participant, err
	}

	return participant, nil
}

func (r *ChallengeRepository) UpdateSubmitChallengeForm(id uint64, updatedStatus dto.UpdateChallengeFormStatusRequest) (*entities.ChallengeFormModels, error) {
	var participant *entities.ChallengeFormModels
	if err := r.db.Model(&entities.ChallengeFormModels{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updatedStatus).Error; err != nil {
		return participant, err
	}
	return participant, nil
}

func (r *ChallengeRepository) GetSubmitChallengeFormByUserAndChallenge(userID uint64) ([]*entities.ChallengeFormModels, error) {
	var submissions []*entities.ChallengeFormModels
	err := r.db.Where("user_id = ?", userID).Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, nil
}
