package service

import (
	"math"

	"github.com/capstone-kelompok-7/backend-disappear/module/article"
	"github.com/capstone-kelompok-7/backend-disappear/module/article/domain"
)

type ArticleRepository struct {
	repo article.RepositoryArticleInterface
}

func NewArticleRepository(repo article.RepositoryArticleInterface) article.ServiceArticleInterface {
	return &ArticleRepository{
		repo: repo,
	}
}

func (s *ArticleRepository) GetAll(page, perPage int) ([]domain.Articles, int64, error) {
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

func (s *ArticleRepository) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
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

func (s *ArticleRepository) GetNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage + 1
	}
	return totalPages
}

func (s *ArticleRepository) GetPrevPage(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}
	return 1
}

func (s *ArticleRepository) GetArticlesByTitle(page, perPage int, title string) ([]domain.Articles, int64, error) {
	articles, err := s.repo.FindByTitle(page, perPage, title)
    if err!= nil {
        return articles, 0, err
    }

    totalItems, err := s.repo.GetTotalArticleCountByTitle(title)
    if err!= nil {
        return articles, 0, err
    }

    return articles, totalItems, nil
}