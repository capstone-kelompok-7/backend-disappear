package service

import (
	"errors"
	"math"
	"time"

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
		Author:  "DISAPPEAR",
	}
	createdArticle, err := s.repo.CreateArticle(value)
	if err != nil {
		return nil, errors.New("Gagal Menambahkan Artikel: " + err.Error())
	}

	return createdArticle, nil
}

func (s *ArticleService) UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error) {
	existingArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan: " + err.Error())
	}

	if existingArticle == nil {
		return nil, errors.New("artikel tidak ditemukan:" + err.Error())
	}

	_, err = s.repo.UpdateArticleById(id, updatedArticle)
	if err != nil {
		return nil, errors.New("Gagal Mengubah Artikel: " + err.Error())
	}

	getUpdatedArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("Gagal mengambil Artikel: " + err.Error())
	}

	return getUpdatedArticle, nil
}

func (s *ArticleService) DeleteArticleById(id uint64) error {
	existingArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return errors.New("Artikel Tidak Ditemukan: " + err.Error())
	}

	if existingArticle == nil {
		return errors.New("Artikel Tidak Ditemukan: " + err.Error())
	}

	err = s.repo.DeleteArticleById(id)
	if err != nil {
		return errors.New("Gagal Menghapus Artikel: " + err.Error())
	}

	return nil
}

func (s *ArticleService) GetAll(page, perPage int) ([]entities.ArticleModels, int64, error) {
	articles, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, 0, errors.New("Artikel Tidak Ditemukan: " + err.Error())
	}

	totalItems, err := s.repo.GetTotalArticleCount()
	if err != nil {
		return nil, 0, errors.New("Gagal Menghitung Total Artikel: " + err.Error())
	}

	return articles, totalItems, nil
}

func (s *ArticleService) GetArticlesByTitle(page, perPage int, title string) ([]entities.ArticleModels, int64, error) {
	articles, err := s.repo.FindByTitle(page, perPage, title)
	if err != nil {
		return nil, 0, errors.New("Artikel Tidak Ditemukan: " + err.Error())
	}

	totalItems, err := s.repo.GetTotalArticleCountByTitle(title)
	if err != nil {
		return nil, 0, errors.New("Gagal Menghitung Total Artikel: " + err.Error())
	}

	return articles, totalItems, nil
}

func (s *ArticleService) GetArticleById(id uint64, incrementViews bool) (*entities.ArticleModels, error) {
	result, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("Artikel Tidak Ditemukan: " + err.Error())
	}

	if incrementViews {
		result.Views++
		if err := s.repo.UpdateArticleViews(result); err != nil {
			return nil, errors.New("Gagal meningkatkan jumlah tayangan artikel: " + err.Error())
		}
	}

	return result, nil
}

func (s *ArticleService) GetArticlesByDateRange(page, perPage int, filterType string) ([]entities.ArticleModels, int64, error) {
	now := time.Now()
	var startDate, endDate time.Time

	switch filterType {
	case "Hari Ini":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endDate = startDate.Add(24 * time.Hour)
	case "Minggu Ini":
		startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
		startDate = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 0, 7)
	case "Bulan Ini":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		nextMonth := startDate.AddDate(0, 1, 0)
		endDate = nextMonth.Add(-time.Second)
	case "Tahun Ini":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		nextYear := startDate.AddDate(1, 0, 0)
		endDate = nextYear.Add(-time.Second)
	default:
		return nil, 0, errors.New("Tipe filter tidak valid")
	}

	result, err := s.repo.GetArticlesByDateRange(page, perPage, startDate, endDate)
	if err != nil {
		return nil, 0, errors.New("Artikel Tidak Ditemukan: " + err.Error())
	}

	totalItems, err := s.repo.GetTotalArticleCountByDateRange(startDate, endDate)
	if err != nil {
		return nil, 0, errors.New("Gagal Menghitung Total Artikel: " + err.Error())
	}

	return result, totalItems, nil
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
