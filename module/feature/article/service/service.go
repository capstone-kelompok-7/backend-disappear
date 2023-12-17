package service

import (
	"errors"
	"math"
	"strings"
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
		return nil, errors.New("gagal menambahkan artikel")
	}

	return createdArticle, nil
}

func (s *ArticleService) UpdateArticleById(id uint64, updatedArticle *entities.ArticleModels) (*entities.ArticleModels, error) {
	existingArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if existingArticle == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	_, err = s.repo.UpdateArticleById(id, updatedArticle)
	if err != nil {
		return nil, errors.New("gagal mengubah artikel")
	}

	getUpdatedArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("gagal mengambil artikel")
	}

	return getUpdatedArticle, nil
}

func (s *ArticleService) DeleteArticleById(id uint64) error {
	existingArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return errors.New("artikel tidak ditemukan")
	}

	err = s.repo.DeleteArticleById(existingArticle.ID)
	if err != nil {
		return errors.New("gagal menghapus artikel")
	}

	return nil
}

func (s *ArticleService) GetAll() ([]*entities.ArticleModels, error) {
	articles, err := s.repo.FindAll()
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	return articles, nil
}

func (s *ArticleService) GetArticlesByTitle(title string) ([]*entities.ArticleModels, error) {
	articles, err := s.repo.FindByTitle(title)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	return articles, nil
}

func (s *ArticleService) GetArticleById(id uint64, incrementViews bool) (*entities.ArticleModels, error) {
	result, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if incrementViews {
		result.Views++
		if err := s.repo.UpdateArticleViews(result); err != nil {
			return nil, errors.New("gagal meningkatkan jumlah tayangan artikel")
		}
	}

	return result, nil
}

func (s *ArticleService) GetArticlesByDateRange(filterType string) ([]*entities.ArticleModels, error) {
	startDate, endDate, err := s.GetFilterDateRange(filterType)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.GetArticlesByDateRange(startDate, endDate)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	return result, nil
}

func (s *ArticleService) BookmarkArticle(bookmark *entities.ArticleBookmarkModels) error {
	articles, err := s.repo.GetArticleById(bookmark.ArticleID)
	if err != nil {
		return errors.New("artikel tidak ditemukan")
	}

	bookmarked, err := s.repo.IsArticleAlreadyBookmarked(bookmark.UserID, articles.ID)
	if err != nil {
		return errors.New("gagal mengecek database")
	}

	if bookmarked {
		return errors.New("artikel telah disimpan")
	}

	newBookmark := &entities.ArticleBookmarkModels{
		UserID:    bookmark.UserID,
		ArticleID: bookmark.ArticleID,
	}
	if err := s.repo.BookmarkArticle(newBookmark); err != nil {
		return errors.New("gagal menyimpan artikel")
	}

	return nil
}

func (s *ArticleService) DeleteBookmarkArticle(userID, articleID uint64) error {
	bookmarked, err := s.repo.IsArticleAlreadyBookmarked(userID, articleID)
	if err != nil {
		return errors.New("gagal mengecek artikel")
	}

	if !bookmarked {
		return errors.New("artikel tidak ditemukan")
	}

	errDelete := s.repo.DeleteBookmarkArticle(userID, articleID)
	if errDelete != nil {
		return errors.New("gagal menghapus artikel tersimpan")
	}
	return nil
}

func (s *ArticleService) GetUserBookmarkArticle(userID uint64) ([]*entities.ArticleBookmarkModels, error) {
	userBookmarkArticle, err := s.repo.GetUserBookmarkArticle(userID)
	if err != nil {
		return nil, errors.New("gagal mendapatkan artikel tersimpan user")
	}
	return userBookmarkArticle, nil
}

func (s *ArticleService) GetLatestArticles() ([]*entities.ArticleModels, error) {
	articles, err := s.repo.GetLatestArticle()
	if err != nil {
		return nil, errors.New("gagal mengambil artikel")
	}

	return articles, nil
}

func (s *ArticleService) GetOldestArticle(page, perPage int) ([]*entities.ArticleModels, int64, error) {
	result, err := s.repo.GetOldestArticle(page, perPage)
	if err != nil {
		return nil, 0, err
	}
	totalItems, err := s.repo.GetTotalArticleCount()
	if err != nil {
		return nil, 0, err
	}
	return result, totalItems, nil
}

func (s *ArticleService) GetAllArticleUser(page, perPage int) ([]*entities.ArticleModels, int64, error) {
	result, err := s.repo.FindAllArticle(page, perPage)
	if err != nil {
		return nil, 0, err
	}
	totalItems, err := s.repo.GetTotalArticleCount()
	if err != nil {
		return nil, 0, err
	}
	return result, totalItems, nil
}

func (s *ArticleService) CalculatePaginationValues(page int, totalItems int, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if page > totalPages {
		page = totalPages
	}

	return page, totalPages
}

func (s *ArticleService) GetNextPage(currentPage int, totalPages int) int {
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

func (s *ArticleService) GetArticlesAlphabet(page, perPage int) ([]*entities.ArticleModels, int64, error) {
	articles, err := s.repo.GetArticleAlphabet(page, perPage)
	if err != nil {
		return nil, 0, errors.New("gagal mengambil artikel")
	}
	totalItems, err := s.repo.GetTotalArticleCount()
	if err != nil {
		return nil, 0, err
	}

	return articles, totalItems, nil
}

func (s *ArticleService) GetArticleMostViews(page, perPage int) ([]*entities.ArticleModels, int64, error) {
	articles, err := s.repo.GetArticleMostViews(page, perPage)
	if err != nil {
		return nil, 0, errors.New("gagal mengambil artikel")
	}
	totalItems, err := s.repo.GetTotalArticleCount()
	if err != nil {
		return nil, 0, err
	}

	return articles, totalItems, nil
}

func (s *ArticleService) GetOtherArticle() ([]*entities.ArticleModels, error) {
	articles, err := s.repo.GetOtherArticle()
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *ArticleService) GetArticleSearchByDateRange(filterType, searchText string) ([]*entities.ArticleModels, error) {
	startDate, endDate, err := s.GetFilterDateRange(filterType)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.SearchArticlesWithDateFilter(searchText, startDate, endDate)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	return result, nil
}

func (s *ArticleService) GetFilterDateRange(filterType string) (time.Time, time.Time, error) {
	filterType = strings.ToLower(filterType)
	now := time.Now()

	switch filterType {
	case "minggu ini":
		startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
		startDate := time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 0, 7)
		return startDate, endDate, nil
	case "bulan ini":
		startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		nextMonth := startDate.AddDate(0, 1, 0)
		endDate := nextMonth.Add(-time.Second)
		return startDate, endDate, nil
	case "tahun ini":
		startDate := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		nextYear := startDate.AddDate(1, 0, 0)
		endDate := nextYear.Add(-time.Second)
		return startDate, endDate, nil
	default:
		return time.Time{}, time.Time{}, errors.New("tipe filter tidak valid")
	}
}
