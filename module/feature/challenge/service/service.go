package service

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge"
	"math"
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
