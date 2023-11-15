package service

import (
	"errors"
	"math"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/article"
)

type ArticleService struct {
	repo article.RepositoryArticleInterface
}

func NewArticleService(repo article.RepositoryArticleInterface) article.ServiceArticleInterface {
	return &ArticleService{
		repo: repo,
	}
}

func (s *ArticleService) CreateArticle(articleData *entities.ArticleModels) (*entities.ArticleModels, error) {
	value := &entities.ArticleModels{
		Title:   articleData.Title,
		Photo:   articleData.Photo,
		Content: articleData.Content,
	}
	createdArticle, err := s.repo.CreateArticle(value)
	if err != nil {
		return nil, err
	}

	return createdArticle, nil
}

func (s *ArticleService) GetAll(page, perPage int) ([]entities.ArticleModels, int64, error) {
	articles, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return articles, 0, err
	}

	totalItems, err := s.repo.GetTotalArticleCount()
	if err != nil {
		return articles, 0, err
	}

	return articles, totalItems, nil
}

func (s *ArticleService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
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

func (s *ArticleService) GetNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}
	return totalPages
}

func (s *ArticleService) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}
	return 1
}

func (s *ArticleService) GetArticlesByTitle(page, perPage int, title string) ([]entities.ArticleModels, int64, error) {
	articles, err := s.repo.FindByTitle(page, perPage, title)
	if err != nil {
		return articles, 0, err
	}

	totalItems, err := s.repo.GetTotalArticleCountByTitle(title)
	if err != nil {
		return articles, 0, err
	}

	return articles, totalItems, nil
}

func (s *ArticleService) UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error) {
	existingArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if existingArticle == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	updatedArticle, err = s.repo.UpdateArticleById(id, updatedArticle)
	if err != nil {
		return nil, errors.New("gagal mengubah artikel ")
	}

	return updatedArticle, nil
}

func (s *ArticleService) GetArticleById(id uint64) (*entities.ArticleModels, error) {
	result, err := s.repo.GetArticleById(id)
	if err!= nil {
        return nil, errors.New("artikel tidak ditemukan")
    }
	return result, nil
}