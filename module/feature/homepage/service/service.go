package service

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/homepage"
)

type HomepageService struct {
	repo homepage.RepositoryHomepageInterface
}

func NewHomepageService(repo homepage.RepositoryHomepageInterface) homepage.ServiceHomepageInterface {
	return &HomepageService{
		repo: repo,
	}
}

func (s *HomepageService) GetBestSellingProducts(limit int) ([]*entities.ProductModels, error) {
	result, err := s.repo.GetBestSellingProducts(limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *HomepageService) GetCategory() ([]*entities.CategoryModels, error) {
	result, err := s.repo.GetFiveCategories()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *HomepageService) GetCarousel() ([]*entities.CarouselModels, error) {
	result, err := s.repo.GetFiveCarousel()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *HomepageService) GetChallenge() ([]*entities.ChallengeModels, error) {
	result, err := s.repo.GetFiveChallenge()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *HomepageService) GetArticle() ([]*entities.ArticleModels, error) {
	result, err := s.repo.GetThreeArticle()
	if err != nil {
		return nil, err
	}
	return result, nil
}
