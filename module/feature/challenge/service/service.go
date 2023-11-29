package service

import (
	"errors"
	"math"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users"
)

type ChallengeService struct {
	repo        challenge.RepositoryChallengeInterface
	userService users.ServiceUserInterface
}

func NewChallengeService(repo challenge.RepositoryChallengeInterface, userService users.ServiceUserInterface) challenge.ServiceChallengeInterface {
	return &ChallengeService{
		repo:        repo,
		userService: userService,
	}
}

func (s *ChallengeService) GetAllChallenges(page, perPage int) ([]*entities.ChallengeModels, int64, error) {
	challenge, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalChallengeCount()
	if err != nil {
		return nil, 0, err
	}

	return challenge, totalItems, nil
}

func (s *ChallengeService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	pageInt := page
	if pageInt <= 0 {
		pageInt = 1
	}

	total_pages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	if pageInt > total_pages {
		pageInt = total_pages
	}

	return pageInt, total_pages
}

func (s *ChallengeService) GetNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}
	return totalPages
}

func (s *ChallengeService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}
	return 1
}

func (s *ChallengeService) GetChallengeByTitle(page, perPage int, title string) ([]*entities.ChallengeModels, int64, error) {
	challenge, err := s.repo.FindByTitle(page, perPage, title)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalChallengeCountByTitle(title)
	if err != nil {
		return nil, 0, err
	}

	return challenge, totalItems, nil
}

func (s *ChallengeService) GetChallengeByStatus(page, perPage int, status string) ([]*entities.ChallengeModels, int64, error) {
	challenges, err := s.repo.FindByStatus(page, perPage, status)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalChallengeCountByStatus(status)
	if err != nil {
		return nil, 0, err
	}

	return challenges, totalItems, nil
}

func (s *ChallengeService) CreateChallenge(newData *entities.ChallengeModels) (*entities.ChallengeModels, error) {
	newChallenge := &entities.ChallengeModels{
		Title:       newData.Title,
		Photo:       newData.Photo,
		StartDate:   newData.StartDate,
		EndDate:     newData.EndDate,
		Description: newData.Description,
		Exp:         newData.Exp,
	}

	currentTime := time.Now()
	if currentTime.After(newChallenge.EndDate) {
		newChallenge.Status = "Kadaluwarsa"
	} else {
		newChallenge.Status = "Belum Kadaluwarsa"
	}

	result, err := s.repo.CreateChallenge(newChallenge)
	if err != nil {
		return result, errors.New("gagal menambahkan tantangan")
	}
	return result, nil
}

func (s *ChallengeService) GetChallengeById(id uint64) (*entities.ChallengeModels, error) {
	result, err := s.repo.GetChallengeById(id)
	if err != nil {
		return result, errors.New("tantangan tidak ditemukan")
	}

	return result, nil
}

func updateIfNotEmpty(target *string, value string) {
	if value != "" {
		*target = value
	}
}

func updateIfNotZero(target *time.Time, value time.Time) {
	if !value.IsZero() {
		*target = value
	}
}

func updateIfNotZeroUint64(target *uint64, value uint64) {
	if value != 0 {
		*target = value
	}
}

func (s *ChallengeService) UpdateChallenge(challengeID uint64, updateData *entities.ChallengeModels) (*entities.ChallengeModels, error) {
	existingChallenge, err := s.repo.GetChallengeById(challengeID)
	if err != nil {
		return &entities.ChallengeModels{}, err
	}

	updateIfNotEmpty(&existingChallenge.Title, updateData.Title)
	updateIfNotEmpty(&existingChallenge.Photo, updateData.Photo)
	updateIfNotZero(&existingChallenge.StartDate, updateData.StartDate)
	updateIfNotZero(&existingChallenge.EndDate, updateData.EndDate)
	updateIfNotEmpty(&existingChallenge.Description, updateData.Description)
	updateIfNotZeroUint64(&existingChallenge.Exp, updateData.Exp)

	currentTime := time.Now()
	if currentTime.After(existingChallenge.EndDate) {
		existingChallenge.Status = "Kadaluwarsa"
	} else {
		existingChallenge.Status = "Belum Kadaluwarsa"
	}

	updateIfNotEmpty(&existingChallenge.Status, updateData.Status)

	updatedChallenge, err := s.repo.UpdateChallenge(challengeID, existingChallenge)
	if err != nil {
		return &entities.ChallengeModels{}, err
	}

	return updatedChallenge, nil
}

func (s *ChallengeService) DeleteChallenge(id uint64) error {
	_, err := s.repo.GetChallengeById(id)
	if err != nil {
		return errors.New("tantangan tidak ditemukan")
	}
	if err := s.repo.DeleteChallenge(id); err != nil {
		return errors.New("gagal menghapus tantangan")
	}
	return nil
}

func (s *ChallengeService) CreateSubmitChallengeForm(form *entities.ChallengeFormModels) (*entities.ChallengeFormModels, error) {
	existingSubmits, err := s.repo.GetSubmitChallengeFormByUserAndChallenge(form.UserID)
	if err != nil {
		return nil, err
	}

	for _, existingSubmit := range existingSubmits {
		if existingSubmit.ChallengeID == form.ChallengeID {
			return nil, errors.New("Anda sudah submit challenge ini sebelumnya")
		}
	}

	challenge, err := s.repo.GetChallengeById(form.ChallengeID)
	if err != nil {
		return nil, err
	}

	if challenge.Status == "Belum Kadaluwarsa" {
		newParticipant := entities.ChallengeFormModels{
			UserID:      form.UserID,
			ChallengeID: form.ChallengeID,
			Username:    form.Username,
			Photo:       form.Photo,
			Status:      "menunggu validasi",
			Exp:         challenge.Exp,
			CreatedAt:   form.CreatedAt,
		}

		result, err := s.repo.CreateSubmitChallengeForm(&newParticipant)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("Tantangan sudah kadaluwarsa, tidak dapat submit")
}

func (s *ChallengeService) GetAllSubmitChallengeForm(page, perPage int) ([]*entities.ChallengeFormModels, int64, error) {
	challenge, err := s.repo.GetAllSubmitChallengeForm(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalSubmitChallengeFormCount()
	if err != nil {
		return nil, 0, err
	}

	return challenge, totalItems, nil
}

func (s *ChallengeService) GetSubmitChallengeFormByStatus(page, perPage int, status string) ([]*entities.ChallengeFormModels, int64, error) {
	challenge, err := s.repo.GetSubmitChallengeFormByStatus(page, perPage, status)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalSubmitChallengeFormCountByStatus(status)
	if err != nil {
		return nil, 0, err
	}

	return challenge, totalItems, nil
}

func (s *ChallengeService) UpdateSubmitChallengeForm(id uint64, updatedData dto.UpdateChallengeFormStatusRequest) (*entities.ChallengeFormModels, error) {
	form, err := s.repo.GetSubmitChallengeFormById(id)
	if err != nil {
		return nil, errors.New("Form tidak ditemukan")
	}

	userID := form.UserID
	user, err := s.userService.GetUsersById(userID)
	if err != nil {
		return nil, errors.New("Gagal mendapatkan data user")
	}

	var changeTotalChallenge int64

	switch {
	case form.Status == "valid" && updatedData.Status == "tidak valid":
		user.Exp -= form.Exp
		changeTotalChallenge = -1
	case form.Status == "tidak valid" && updatedData.Status == "valid":
		user.Exp += form.Exp
		changeTotalChallenge = 1
	case form.Status == "menunggu validasi" && updatedData.Status == "valid":
		user.Exp += form.Exp
		changeTotalChallenge = 1
	case form.Status == "valid" && updatedData.Status == "menunggu validasi":
		user.Exp -= form.Exp
		changeTotalChallenge = -1
	}

	user.TotalChallenge += uint64(changeTotalChallenge)

	_, err = s.userService.UpdateUserChallengeFollow(userID, user.TotalChallenge)
	if err != nil {
		return nil, errors.New("Gagal menyimpan perubahan total challenge user ke database")
	}

	_, err = s.userService.UpdateUserExp(userID, user.Exp)
	if err != nil {
		return nil, errors.New("Gagal menyimpan perubahan exp user ke database")
	}

	result, err := s.repo.UpdateSubmitChallengeForm(id, updatedData)
	if err != nil {
		return nil, errors.New("Gagal memperbarui form")
	}

	return result, nil
}

func (s *ChallengeService) GetSubmitChallengeFormById(id uint64) (*entities.ChallengeFormModels, error) {
	result, err := s.repo.GetSubmitChallengeFormById(id)
	if err != nil {
		return result, errors.New("form tidak ditemukan")
	}

	return result, nil
}
