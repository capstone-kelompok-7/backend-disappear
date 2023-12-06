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
		return nil, errors.New("gagal menambahkan artikel: " + err.Error())
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
		return nil, errors.New("gagal mengubah artikel: " + err.Error())
	}

	getUpdatedArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("gagal mengambil artikel: " + err.Error())
	}

	return getUpdatedArticle, nil
}

func (s *ArticleService) DeleteArticleById(id uint64) error {
	existingArticle, err := s.repo.GetArticleById(id)
	if err != nil {
		return errors.New("artikel tidak ditemukan: " + err.Error())
	}

	if existingArticle == nil {
		return errors.New("artikel tidak ditemukan: " + err.Error())
	}

	err = s.repo.DeleteArticleById(id)
	if err != nil {
		return errors.New("gagal menghapus artikel: " + err.Error())
	}

	return nil
}

func (s *ArticleService) GetAll() ([]*entities.ArticleModels, error) {
	articles, err := s.repo.FindAll()
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan: " + err.Error())
	}

	return articles, nil
}

func (s *ArticleService) GetArticlesByTitle(title string) ([]*entities.ArticleModels, error) {
	articles, err := s.repo.FindByTitle(title)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan: " + err.Error())
	}

	return articles, nil
}

func (s *ArticleService) GetArticleById(id uint64, incrementViews bool) (*entities.ArticleModels, error) {
	result, err := s.repo.GetArticleById(id)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan: " + err.Error())
	}

	if incrementViews {
		result.Views++
		if err := s.repo.UpdateArticleViews(result); err != nil {
			return nil, errors.New("gagal meningkatkan jumlah tayangan artikel: " + err.Error())
		}
	}

	return result, nil
}

func (s *ArticleService) GetArticlesByDateRange(filterType string) ([]*entities.ArticleModels, error) {
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
		return nil, errors.New("tipe filter tidak valid")
	}

	result, err := s.repo.GetArticlesByDateRange(startDate, endDate)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan: " + err.Error())
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
		return nil, errors.New("gagal mengambil artikel: " + err.Error())
	}

	return articles, nil
}

func (s *ArticleService) GetArticlesByViews(sortType string) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	var err error
	if sortType == "asc" {
		articles, err = s.repo.GetArticleByViewsAsc()
		if err != nil {
			return nil, errors.New("gagal mengambil artikel: " + err.Error())
		}
	} else if sortType == "desc" {
		articles, err = s.repo.GetArticleByViewsDesc()
		if err != nil {
			return nil, errors.New("gagal mengambil artikel: " + err.Error())
		}
	}

	return articles, nil
}

func (s *ArticleService) GetArticlesBySortedTitle(sortType string) ([]*entities.ArticleModels, error) {
	var articles []*entities.ArticleModels
	var err error
	if sortType == "asc" {
		articles, err = s.repo.GetArticleByTitleAsc()
		if err != nil {
			return nil, errors.New("gagal mengambil artikel: " + err.Error())
		}
		return articles, nil
	} else if sortType == "desc" {
		articles, err = s.repo.GetArticleByTitleDesc()
		if err != nil {
			return nil, errors.New("gagal mengambil artikel: " + err.Error())
		}

		return articles, nil
	}

	return nil, nil
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

func (s *ArticleService) GetArticlePreferences(userID uint64, page, perPage int) ([]*entities.ArticleModels, int64, error) {
	result, err := s.repo.FindAllByUserPreference(userID, page, perPage)
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
		return nil, 0, errors.New("gagal mengambil artikel: " + err.Error())
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
		return nil, 0, errors.New("gagal mengambil artikel: " + err.Error())
	}
	totalItems, err := s.repo.GetTotalArticleCount()
	if err != nil {
		return nil, 0, err
	}

	return articles, totalItems, nil
}
