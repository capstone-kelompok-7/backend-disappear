package service

import (
	"errors"
	"math"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
)

type ChallengeService struct {
	repo challenge.RepositoryChallengeInterface
}

func NewChallengeService(repo challenge.RepositoryChallengeInterface) challenge.ServiceChallengeInterface {
	return &ChallengeService{
		repo: repo,
	}
}

func (s *ChallengeService) GetAllChallenges(page, perPage int) ([]entities.ChallengeModels, int64, error) {
	challenge, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return challenge, 0, err
	}

	totalItems, err := s.repo.GetTotalChallengeCount()
	if err != nil {
		return challenge, 0, err
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

func (s *ChallengeService) GetChallengeByTitle(page, perPage int, title string) ([]entities.ChallengeModels, int64, error) {
	challenge, err := s.repo.FindByTitle(page, perPage, title)
	if err != nil {
		return challenge, 0, err
	}

	totalItems, err := s.repo.GetTotalChallengeCountByTitle(title)
	if err != nil {
		return challenge, 0, err
	}

	return challenge, totalItems, nil
}

func (s *ChallengeService) CreateChallenge(newData entities.ChallengeModels) (entities.ChallengeModels, error) {
	newChallenge := entities.ChallengeModels{
		Title:       newData.Title,
		Photo:       newData.Photo,
		StartDate:   newData.StartDate,
		EndDate:     newData.EndDate,
		Description: newData.Description,
		Exp:         newData.Exp,
	}

	currentTime := time.Now()
	if currentTime.After(newChallenge.EndDate) {
		newChallenge.Status = "Berakhir"
	} else {
		newChallenge.Status = "Berlangsung"
	}

	result, err := s.repo.CreateChallenge(newChallenge)
	if err != nil {
		return result, errors.New("gagal menambahkan tantangan")
	}
	return result, nil
}

func (s *ChallengeService) GetChallengeById(id uint64) (entities.ChallengeModels, error) {
	result, err := s.repo.GetChallengeById(id)
	if err != nil {
		return result, errors.New("challenge tidak ditemukan")
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

func (s *ChallengeService) UpdateChallenge(challengeID uint64, updateData entities.ChallengeModels) (entities.ChallengeModels, error) {
	existingChallenge, err := s.repo.GetChallengeById(challengeID)
	if err != nil {
		return entities.ChallengeModels{}, err
	}

	updateIfNotEmpty(&existingChallenge.Title, updateData.Title)
	updateIfNotEmpty(&existingChallenge.Photo, updateData.Photo)
	updateIfNotZero(&existingChallenge.StartDate, updateData.StartDate)
	updateIfNotZero(&existingChallenge.EndDate, updateData.EndDate)
	updateIfNotEmpty(&existingChallenge.Description, updateData.Description)
	updateIfNotZeroUint64(&existingChallenge.Exp, updateData.Exp)

	currentTime := time.Now()
	if currentTime.After(existingChallenge.EndDate) {
		existingChallenge.Status = "Berakhir"
	} else {
		existingChallenge.Status = "Berlangsung"
	}

	updateIfNotEmpty(&existingChallenge.Status, updateData.Status)

	updatedChallenge, err := s.repo.UpdateChallenge(challengeID, existingChallenge)
	if err != nil {
		return entities.ChallengeModels{}, err
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
