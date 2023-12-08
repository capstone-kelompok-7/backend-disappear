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
	err := r.db.Model(&entities.ChallengeModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) FindByTitle(page, perpage int, title string) ([]*entities.ChallengeModels, error) {
	var challenge []*entities.ChallengeModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("title LIKE? AND deleted_at IS NULL", "%"+title+"%").Find(&challenge).Error
	if err != nil {
		return challenge, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) GetTotalChallengeCountByTitle(title string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeModels{}).Where("title LIKE? AND deleted_at IS NULL", "%"+title+"%").Count(&count).Error
	return count, err
}

// func (r *CarouselRepository) GetTotalCarouselCountByName(name string) (int64, error) {
// 	var count int64
// 	query := r.db.Model(&entities.CarouselModels{}).Where("deleted_at IS NULL")

// 	if name != "" {
// 		query = query.Where("name LIKE ?", "%"+name+"%")
// 	}

// 	err := query.Count(&count).Error
// 	return count, err
// }

func (r *ChallengeRepository) FindByStatus(page, perpage int, status string) ([]*entities.ChallengeModels, error) {
	var challenge []*entities.ChallengeModels
	offset := (page - 1) * perpage
	err := r.db.Offset(offset).Limit(perpage).Where("status = ? AND deleted_at IS NULL", status).Find(&challenge).Error
	if err != nil {
		return challenge, err
	}
	return challenge, nil
}

func (r *ChallengeRepository) GetTotalChallengeCountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeModels{}).Where("status = ? AND deleted_at IS NULL", status).Count(&count).Error
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
	if err := r.db.Model(&entities.ChallengeModels{}).Where("id = ? AND deleted_at IS NULL", challengeID).Updates(updatedChallenge).Error; err != nil {
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
	err := r.db.Where("status = ? AND deleted_at IS NULL", status).Offset(offset).Limit(perpage).Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *ChallengeRepository) GetTotalSubmitChallengeFormCount() (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeFormModels{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) GetTotalSubmitChallengeFormCountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeFormModels{}).Where("status = ? AND deleted_at IS NULL", status).Count(&count).Error
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
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *ChallengeRepository) GetSubmitChallengeFormByDateRange(page, perpage int, startDate, endDate time.Time) ([]*entities.ChallengeFormModels, error) {
	var participant []*entities.ChallengeFormModels
	offset := (page - 1) * perpage
	if err := r.db.Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).Offset(offset).Limit(perpage).Find(&participant).Error; err != nil {
		return nil, err
	}
	return participant, nil
}

func (r *ChallengeRepository) GetTotalSubmitChallengeFormCountByDateRange(startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeFormModels{}).Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).Count(&count).Error
	return count, err
}

func (r *ChallengeRepository) GetSubmitChallengeFormByStatusAndDate(page, perPage int, filterStatus string, startDate, endDate time.Time) ([]*entities.ChallengeFormModels, error) {
	var participants []*entities.ChallengeFormModels
	offset := (page - 1) * perPage

	err := r.db.Where("status = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL", filterStatus, startDate, endDate).
		Offset(offset).
		Limit(perPage).
		Find(&participants).Error

	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (r *ChallengeRepository) GetTotalSubmitChallengeFormCountByStatusAndDate(filterStatus string, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&entities.ChallengeFormModels{}).
		Where("status = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL", filterStatus, startDate, endDate).
		Count(&count).
		Error

	return count, err
}

func (r *ChallengeRepository) GetChallengesBySearchAndStatus(page, perPage int, search, status string) ([]*entities.ChallengeModels, int64, error) {
	var challenges []*entities.ChallengeModels
	offset := (page - 1) * perPage

	db := r.db.Model(&entities.ChallengeModels{}).Where("deleted_at IS NULL")

	if search != "" {
		db = db.Where("title LIKE ? ", "%"+search+"%")
	}

	if status != "" {
		db = db.Where("status = ? ", status)
	}

	err := db.Offset(offset).Limit(perPage).Find(&challenges).Error
	if err != nil {
		return nil, 0, err
	}

	var totalItems int64
	if err := db.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	return challenges, totalItems, nil
}
