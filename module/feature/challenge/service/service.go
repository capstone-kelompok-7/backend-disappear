package service

import (
	"errors"
	"math"
	"strings"
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
			return nil, errors.New("anda sudah submit tantangan ini sebelumnya")
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

	return nil, errors.New("tantangan sudah kadaluwarsa, tidak dapat submit")
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
		return nil, errors.New("formulir tidak ditemukan")
	}

	user, err := s.userService.GetUsersById(form.UserID)
	if err != nil {
		return nil, errors.New("pengguna tidak ada")
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

	result, err := s.repo.UpdateSubmitChallengeForm(id, updatedData)
	if err != nil {
		return nil, errors.New("gagal memperbarui formulir")
	}

	_, err = s.userService.UpdateUserExp(user.ID, user.Exp)
	if err != nil {
		return nil, errors.New("gagal menyimpan perubahan exp user ke database")
	}

	_, err = s.userService.UpdateUserChallengeFollow(user.ID, user.TotalChallenge)
	if err != nil {
		return nil, errors.New("gagal menyimpan perubahan total tantangan user ke database")
	}

	return result, nil
}

func (s *ChallengeService) GetSubmitChallengeFormById(id uint64) (*entities.ChallengeFormModels, error) {
	result, err := s.repo.GetSubmitChallengeFormById(id)
	if err != nil {
		return result, errors.New("formulir tidak ditemukan")
	}

	return result, nil
}

func getDatesFromFilterType(filterType string) (time.Time, time.Time, error) {
	filterType = strings.ToLower(filterType)
	now := time.Now()
	var startDate, endDate time.Time

	switch filterType {
	case "hari ini":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endDate = startDate.Add(24 * time.Hour)
	case "minggu ini":
		startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
		startDate = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 0, 7)
	case "bulan ini":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		nextMonth := startDate.AddDate(0, 1, 0)
		endDate = nextMonth.Add(-time.Second)
	case "tahun ini":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		nextYear := startDate.AddDate(1, 0, 0)
		endDate = nextYear.Add(-time.Second)
	default:
		return time.Time{}, time.Time{}, errors.New("tipe filter tidak valid")
	}

	return startDate, endDate, nil
}

func (s *ChallengeService) GetSubmitChallengeFormByDateRange(page, perPage int, filterType string) ([]*entities.ChallengeFormModels, int64, error) {
	startDate, endDate, err := getDatesFromFilterType(filterType)
	if err != nil {
		return nil, 0, err
	}

	result, err := s.repo.GetSubmitChallengeFormByDateRange(page, perPage, startDate, endDate)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan form pengumpulan: " + err.Error())
	}

	totalItems, err := s.repo.GetTotalSubmitChallengeFormCountByDateRange(startDate, endDate)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total form pengumpulan: " + err.Error())
	}

	return result, totalItems, nil
}

func (s *ChallengeService) GetSubmitChallengeFormByStatusAndDate(page, perPage int, filterStatus string, filterType string) ([]*entities.ChallengeFormModels, int64, error) {
	startDate, endDate, err := getDatesFromFilterType(filterType)
	if err != nil {
		return nil, 0, err
	}

	participants, err := s.repo.GetSubmitChallengeFormByStatusAndDate(page, perPage, filterStatus, startDate, endDate)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan form pengumpulan: " + err.Error())
	}

	totalItems, err := s.repo.GetTotalSubmitChallengeFormCountByStatusAndDate(filterStatus, startDate, endDate)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan total form pengumpulan: " + err.Error())
	}

	return participants, totalItems, nil
}

func (s *ChallengeService) GetChallengesBySearchAndStatus(page, perPage int, search, status string) ([]*entities.ChallengeModels, int64, error) {
	challenges, totalItems, err := s.repo.GetChallengesBySearchAndStatus(page, perPage, search, status)
	if err != nil {
		return nil, 0, errors.New("gagal mendapatkan tantangan: " + err.Error())
	}

	return challenges, totalItems, nil
}
