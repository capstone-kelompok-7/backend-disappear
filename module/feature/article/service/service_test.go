package service

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/article/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleService_CalculatePaginationValues(t *testing.T) {
	service := &ArticleService{}

	t.Run("Page less than or equal to zero should default to 1", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(0, 100, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page exceeds total pages should set to total pages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(15, 100, 8)

		assert.Equal(t, 13, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page within limits should return correct values", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(2, 100, 8)

		assert.Equal(t, 2, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Total items not perfectly divisible by perPage should round totalPages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(1, 95, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 12, totalPages)
	})
}

func TestArticleService_GetNextPage(t *testing.T) {
	service := &ArticleService{}

	t.Run("Next Page Within Total Pages", func(t *testing.T) {
		currentPage := 3
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, currentPage+1, nextPage)
	})

	t.Run("Next Page Equal to Total Pages", func(t *testing.T) {
		currentPage := 5
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, totalPages, nextPage)
	})
}

func TestArticleService_GetPrevPage(t *testing.T) {
	service := &ArticleService{}

	t.Run("Previous Page Within Bounds", func(t *testing.T) {
		currentPage := 3

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage-1, prevPage)
	})

	t.Run("Previous Page at Lower Bound", func(t *testing.T) {
		currentPage := 1

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage, prevPage)
	})
}

func TestArticleService_GetAll(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	articles := []*entities.ArticleModels{
		{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
	}

	t.Run("Success Case - Articles Found", func(t *testing.T) {
		repo.On("FindAll").Return(articles, nil).Once()

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, articles, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Articles Not Found", func(t *testing.T) {
		expectedErr := errors.New("artikel tidak ditemukan")

		repo.On("FindAll").Return(nil, expectedErr).Once()

		articles, err := service.GetAll()

		assert.Error(t, err)
		assert.Nil(t, articles)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetArticleByTitle(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	articles := []*entities.ArticleModels{
		{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
	}

	name := "test"

	t.Run("Success Case - Articles Found by Title", func(t *testing.T) {
		repo.On("FindByTitle", name).Return(articles, nil).Once()

		result, err := service.GetArticlesByTitle(name)

		assert.NoError(t, err)
		assert.Equal(t, articles, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Articles by Title", func(t *testing.T) {
		expectedErr := errors.New("artikel tidak ditemukan")

		repo.On("FindByTitle", name).Return(nil, expectedErr).Once()

		result, err := service.GetArticlesByTitle(name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetArticleById(t *testing.T) {
	repo := new(mocks.RepositoryArticleInterface)
	service := NewArticleService(repo)

	article := &entities.ArticleModels{
		ID:        1,
		Title:     "Article 1",
		Photo:     "article2.jpg",
		Content:   "this is article 1",
		Author:    "admin",
		CreatedAt: time.Now(),
		Views:     1,
	}

	expectedArticle := &entities.ArticleModels{
		ID:        article.ID,
		Title:     article.Title,
		Photo:     article.Photo,
		Content:   article.Content,
		Author:    article.Author,
		CreatedAt: article.CreatedAt,
		Views:     article.Views,
	}

	t.Run("Success Case - Article With Specific Id Found", func(t *testing.T) {
		articleID := uint64(1)
		incrementViews := true

		repo.On("GetArticleById", articleID).Return(expectedArticle, nil).Once()
		repo.On("UpdateArticleViews", expectedArticle).Return(nil).Once()

		result, err := service.GetArticleById(articleID, incrementViews)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedArticle.ID, result.ID)
		assert.Equal(t, expectedArticle.Title, result.Title)
		assert.Equal(t, expectedArticle.Photo, result.Photo)
		assert.Equal(t, expectedArticle.Content, result.Content)
		assert.Equal(t, expectedArticle.Author, result.Author)
		assert.Equal(t, expectedArticle.CreatedAt, result.CreatedAt)
		assert.Equal(t, expectedArticle.Views, result.Views)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article With Specific Id Not Found", func(t *testing.T) {
		articleID := uint64(2)
		incrementViews := true
		expectedErr := errors.New("artikel tidak ditemukan")

		repo.On("GetArticleById", articleID).Return(nil, expectedErr).Once()

		result, err := service.GetArticleById(articleID, incrementViews)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Increase the Number of Article Views", func(t *testing.T) {
		articleID := uint64(3)
		incrementViews := true
		expectedErr := errors.New("gagal meningkatkan jumlah tayangan artikel")

		repo.On("GetArticleById", articleID).Return(expectedArticle, nil).Once()
		repo.On("UpdateArticleViews", expectedArticle).Return(expectedErr).Once()

		result, err := service.GetArticleById(articleID, incrementViews)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_CreateArticle(t *testing.T) {
	repo := new(mocks.RepositoryArticleInterface)
	service := NewArticleService(repo)

	articleData := &entities.ArticleModels{
		Title:   "New Article",
		Photo:   "new_article.jpg",
		Content: "This is a new article",
	}

	expectedCreatedArticle := &entities.ArticleModels{
		ID:        1,
		Title:     articleData.Title,
		Photo:     articleData.Photo,
		Content:   articleData.Content,
		Author:    "DISAPPEAR",
		CreatedAt: time.Now(),
		Views:     0,
	}

	t.Run("Success Case - Article Created", func(t *testing.T) {
		repo.On("CreateArticle", mock.AnythingOfType("*entities.ArticleModels")).Return(expectedCreatedArticle, nil).Once()

		result, err := service.CreateArticle(articleData)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCreatedArticle.ID, result.ID)
		assert.Equal(t, expectedCreatedArticle.Title, result.Title)
		assert.Equal(t, expectedCreatedArticle.Photo, result.Photo)
		assert.Equal(t, expectedCreatedArticle.Content, result.Content)
		assert.Equal(t, expectedCreatedArticle.Author, result.Author)
		assert.Equal(t, expectedCreatedArticle.CreatedAt, result.CreatedAt)
		assert.Equal(t, expectedCreatedArticle.Views, result.Views)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Create Article", func(t *testing.T) {
		expectedErr := errors.New("gagal menambahkan artikel")
		repo.On("CreateArticle", mock.AnythingOfType("*entities.ArticleModels")).Return(nil, expectedErr).Once()

		result, err := service.CreateArticle(articleData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_UpdateArticleById(t *testing.T) {
	repo := new(mocks.RepositoryArticleInterface)
	service := NewArticleService(repo)

	existingArticle := &entities.ArticleModels{
		ID:        1,
		Title:     "Existing Article",
		Photo:     "existing_article.jpg",
		Content:   "This is an existing article",
		Author:    "DISAPPEAR",
		CreatedAt: time.Now(),
		Views:     0,
	}

	updatedArticle := &entities.ArticleModels{
		Title:   "Updated Article",
		Photo:   "updated_article.jpg",
		Content: "This is an updated article",
	}

	t.Run("Success Case - Success Update Article", func(t *testing.T) {
		articleID := uint64(1)

		repo.On("GetArticleById", articleID).Return(existingArticle, nil).Once()
		repo.On("UpdateArticleById", articleID, updatedArticle).Return(nil, nil).Once()
		repo.On("GetArticleById", articleID).Return(updatedArticle, nil).Once()

		result, err := service.UpdateArticleById(articleID, updatedArticle)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedArticle.Title, result.Title)
		assert.Equal(t, updatedArticle.Photo, result.Photo)
		assert.Equal(t, updatedArticle.Content, result.Content)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article With Specific Id Not Found", func(t *testing.T) {
		articleID := uint64(2)
		expectedErr := errors.New("artikel tidak ditemukan")

		repo.On("GetArticleById", articleID).Return(nil, expectedErr).Once()

		result, err := service.UpdateArticleById(articleID, updatedArticle)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article With Specific Id Found But Update Failed", func(t *testing.T) {
		articleID := uint64(3)
		expectedErr := errors.New("gagal mengubah artikel")

		repo.On("GetArticleById", articleID).Return(existingArticle, nil).Once()
		repo.On("UpdateArticleById", articleID, updatedArticle).Return(nil, expectedErr).Once()

		result, err := service.UpdateArticleById(articleID, updatedArticle)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Find Updated Article", func(t *testing.T) {
		articleID := uint64(4)
		expectedErr := errors.New("gagal mengambil artikel")

		repo.On("GetArticleById", articleID).Return(existingArticle, nil).Once()
		repo.On("UpdateArticleById", articleID, updatedArticle).Return(nil, nil).Once()
		repo.On("GetArticleById", articleID).Return(nil, expectedErr).Once()

		result, err := service.UpdateArticleById(articleID, updatedArticle)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article Is Nil", func(t *testing.T) {
		articleID := uint64(5)
		expectedErr := errors.New("artikel tidak ditemukan")

		repo.On("GetArticleById", articleID).Return(nil, nil).Once()

		result, err := service.UpdateArticleById(articleID, updatedArticle)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

}

func TestArticleService_DeleteArticleById(t *testing.T) {
	repo := new(mocks.RepositoryArticleInterface)
	service := NewArticleService(repo)

	existingArticle := &entities.ArticleModels{
		ID:        1,
		Title:     "Existing Article",
		Photo:     "existing_article.jpg",
		Content:   "This is an existing article",
		Author:    "DISAPPEAR",
		CreatedAt: time.Now(),
		Views:     0,
	}

	t.Run("Success Case - Successful Delete Article", func(t *testing.T) {
		articleID := uint64(1)

		repo.On("GetArticleById", articleID).Return(existingArticle, nil).Once()
		repo.On("DeleteArticleById", articleID).Return(nil).Once()

		err := service.DeleteArticleById(articleID)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article With Specific Id Not Found", func(t *testing.T) {
		articleID := uint64(2)
		expectedErr := errors.New("artikel tidak ditemukan")

		repo.On("GetArticleById", articleID).Return(nil, expectedErr).Once()

		err := service.DeleteArticleById(articleID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article With Specific Id Found But Delete Failed", func(t *testing.T) {
		repo.On("GetArticleById", uint64(1)).Return(&entities.ArticleModels{ID: 1}, nil)
		repo.On("DeleteArticleById", uint64(1)).Return(errors.New("gagal menghapus artikel")).Once()

		err := service.DeleteArticleById(1)

		expectedErr := errors.New("gagal menghapus artikel")

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetArticlesByDateRange(t *testing.T) {
	repo := new(mocks.RepositoryArticleInterface)
	service := NewArticleService(repo)

	t.Run("Success Case - Success Get Articles By Date Range", func(t *testing.T) {
		filterType := "bulan ini"
		startDate, endDate, _ := service.GetFilterDateRange(filterType)

		expectedArticles := []*entities.ArticleModels{
			{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: startDate.Add(time.Hour), Views: 1},
			{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: startDate.Add(-time.Hour), Views: 1},
		}

		repo.On("GetArticlesByDateRange", startDate, endDate).Return(expectedArticles, nil).Once()

		result, err := service.GetArticlesByDateRange(filterType)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedArticles, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get Date Range", func(t *testing.T) {
		filterType := "invalid type"
		result, err := service.GetArticlesByDateRange(filterType)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "tipe filter tidak valid", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article Not Found", func(t *testing.T) {
		filterType := "bulan ini"
		startDate, endDate, _ := service.GetFilterDateRange(filterType)

		repo.On("GetArticlesByDateRange", startDate, endDate).Return(nil, errors.New("artikel tidak ditemukan")).Once()

		result, err := service.GetArticlesByDateRange(filterType)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "artikel tidak ditemukan", err.Error())
		repo.AssertExpectations(t)
	})

}

func TestArticleService_BookmarkArticle(t *testing.T) {
	repo := new(mocks.RepositoryArticleInterface)
	service := NewArticleService(repo)

	userID := uint64(1)
	articleID := uint64(2)

	t.Run("Success Case - Success Bookmark Article", func(t *testing.T) {
		bookmark := &entities.ArticleBookmarkModels{
			UserID:    userID,
			ArticleID: articleID,
		}

		repo.On("GetArticleById", articleID).Return(&entities.ArticleModels{ID: articleID}, nil).Once()
		repo.On("IsArticleAlreadyBookmarked", userID, articleID).Return(false, nil).Once()
		repo.On("BookmarkArticle", bookmark).Return(nil).Once()

		err := service.BookmarkArticle(bookmark)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article Not Found", func(t *testing.T) {
		bookmark := &entities.ArticleBookmarkModels{
			UserID:    userID,
			ArticleID: articleID,
		}

		repo.On("GetArticleById", articleID).Return(nil, errors.New("artikel tidak ditemukan")).Once()

		err := service.BookmarkArticle(bookmark)

		assert.Error(t, err)
		assert.Equal(t, "artikel tidak ditemukan", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article Found But Failed Check Bookmark Database", func(t *testing.T) {
		bookmark := &entities.ArticleBookmarkModels{
			UserID:    userID,
			ArticleID: articleID,
		}

		repo.On("GetArticleById", articleID).Return(&entities.ArticleModels{ID: articleID}, nil).Once()
		repo.On("IsArticleAlreadyBookmarked", userID, articleID).Return(false, errors.New("gagal mengecek database")).Once()

		err := service.BookmarkArticle(bookmark)

		assert.Error(t, err)
		assert.Equal(t, "gagal mengecek database", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Bookmark Article", func(t *testing.T) {
		bookmark := &entities.ArticleBookmarkModels{
			UserID:    userID,
			ArticleID: articleID,
		}

		repo.On("GetArticleById", articleID).Return(&entities.ArticleModels{ID: articleID}, nil).Once()
		repo.On("IsArticleAlreadyBookmarked", userID, articleID).Return(false, nil).Once()
		repo.On("BookmarkArticle", bookmark).Return(errors.New("gagal menyimpan artikel")).Once()

		err := service.BookmarkArticle(bookmark)

		assert.Error(t, err)
		assert.Equal(t, "gagal menyimpan artikel", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article Already Bookmarked", func(t *testing.T) {
		bookmark := &entities.ArticleBookmarkModels{
			UserID:    userID,
			ArticleID: articleID,
		}

		repo.On("GetArticleById", articleID).Return(&entities.ArticleModels{ID: articleID}, nil).Once()
		repo.On("IsArticleAlreadyBookmarked", userID, articleID).Return(true, nil).Once()

		err := service.BookmarkArticle(bookmark)

		assert.Error(t, err)
		assert.Equal(t, "artikel telah disimpan", err.Error())
		repo.AssertExpectations(t)
	})

}

func TestArticleService_GetLatestArticles(t *testing.T) {
	repo := new(mocks.RepositoryArticleInterface)
	service := NewArticleService(repo)

	t.Run("Success Case - Success Get Latest Articles", func(t *testing.T) {
		expectedArticles := []*entities.ArticleModels{
			{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
			{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		}

		repo.On("GetLatestArticle").Return(expectedArticles, nil).Once()

		result, err := service.GetLatestArticles()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedArticles, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get Latest Articles", func(t *testing.T) {
		repo.On("GetLatestArticle").Return(nil, errors.New("gagal mengambil artikel")).Once()

		result, err := service.GetLatestArticles()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "gagal mengambil artikel", err.Error())
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetOldestArticle(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	articles := []*entities.ArticleModels{
		{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
	}

	t.Run("Success Case - Success Get Oldest Articles", func(t *testing.T) {
		expectedTotalItems := int64(10)

		repo.On("GetOldestArticle", 1, 10).Return(articles, nil).Once()
		repo.On("GetTotalArticleCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetOldestArticle(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(articles), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get Oldest Articles", func(t *testing.T) {
		expectedErr := errors.New("failed to get oldest article")

		repo.On("GetOldestArticle", 1, 10).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetOldestArticle(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Article", func(t *testing.T) {
		expectedErr := errors.New("failed to get total article count")

		repo.On("GetOldestArticle", 1, 10).Return(articles, nil).Once()
		repo.On("GetTotalArticleCount").Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetOldestArticle(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetArticlesAlphabet(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	articles := []*entities.ArticleModels{
		{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
	}

	t.Run("Success Case - Success Get Articles by Alphabet", func(t *testing.T) {
		expectedTotalItems := int64(10)

		repo.On("GetArticleAlphabet", 1, 10).Return(articles, nil).Once()
		repo.On("GetTotalArticleCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetArticlesAlphabet(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(articles), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get Articles by Alphabet", func(t *testing.T) {
		expectedErr := errors.New("gagal mengambil artikel")

		repo.On("GetArticleAlphabet", 1, 10).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetArticlesAlphabet(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Article", func(t *testing.T) {
		expectedErr := errors.New("failed to get total article count")

		repo.On("GetArticleAlphabet", 1, 10).Return(articles, nil).Once()
		repo.On("GetTotalArticleCount").Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetArticlesAlphabet(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetArticleMostViews(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	articles := []*entities.ArticleModels{
		{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
	}

	t.Run("Success Case - Success Get Most View Articles", func(t *testing.T) {
		expectedTotalItems := int64(10)

		repo.On("GetArticleMostViews", 1, 10).Return(articles, nil).Once()
		repo.On("GetTotalArticleCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetArticleMostViews(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(articles), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get Most View Articles", func(t *testing.T) {
		expectedErr := errors.New("gagal mengambil artikel")
		repo.On("GetArticleMostViews", 1, 10).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetArticleMostViews(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Article", func(t *testing.T) {
		expectedErr := errors.New("failed to get total article count")

		repo.On("GetArticleMostViews", 1, 10).Return(articles, nil).Once()
		repo.On("GetTotalArticleCount").Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetArticleMostViews(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetOtherArticle(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	articles := []*entities.ArticleModels{
		{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
	}

	t.Run("Success Case - Success Get Other Articles", func(t *testing.T) {
		repo.On("GetOtherArticle").Return(articles, nil).Once()

		result, err := service.GetOtherArticle()

		assert.NoError(t, err)
		assert.Equal(t, len(articles), len(result))
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get Other Articles", func(t *testing.T) {
		expectedErr := errors.New("failed to get other articles")

		repo.On("GetOtherArticle").Return(nil, expectedErr).Once()

		result, err := service.GetOtherArticle()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetArticleSearchByDateRange(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	articles := []*entities.ArticleModels{
		{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
	}

	t.Run("Success Case - Success Get Article Search By Date Range", func(t *testing.T) {
		filterType := "bulan ini"
		expectedStartDate := time.Date(2023, time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
		expectedEndDate := time.Date(2023, time.Now().Month(), 31, 23, 59, 59, 0, time.UTC)

		repo.On("SearchArticlesWithDateFilter", "", expectedStartDate, expectedEndDate).Return(articles, nil).Once()

		result, err := service.GetArticleSearchByDateRange(filterType, "")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, articles, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get Filter Date Range", func(t *testing.T) {
		filterType := "invalid type"
		expectedErr := errors.New("tipe filter tidak valid")

		result, err := service.GetArticleSearchByDateRange(filterType, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - No Articles Found in Date Range", func(t *testing.T) {
		filterType := "bulan ini"
		expectedStartDate := time.Date(2023, time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
		expectedEndDate := time.Date(2023, time.Now().Month(), 31, 23, 59, 59, 0, time.UTC)
		expectedErr := errors.New("artikel tidak ditemukan")
		repo.On("SearchArticlesWithDateFilter", "", expectedStartDate, expectedEndDate).Return(nil, expectedErr).Once()

		result, err := service.GetArticleSearchByDateRange(filterType, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})

}

func TestArticleService_GetAllArticleUser(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	articles := []*entities.ArticleModels{
		{ID: 1, Title: "Article 1", Photo: "article1.jpg", Content: "this is content of article 1", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
		{ID: 2, Title: "Article 2", Photo: "article2.jpg", Content: "this is content of article 2", Author: "DISAPPEAR", CreatedAt: time.Now(), Views: 1},
	}

	t.Run("Success Case - Success Get All User Articles", func(t *testing.T) {
		expectedTotalItems := int64(10)

		repo.On("FindAllArticle", 1, 10).Return(articles, nil).Once()
		repo.On("GetTotalArticleCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetAllArticleUser(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(articles), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get All User Articles", func(t *testing.T) {
		expectedErr := errors.New("failed to get user articles")

		repo.On("FindAllArticle", 1, 10).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetAllArticleUser(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Article", func(t *testing.T) {
		expectedErr := errors.New("failed to get total article count")
		repo.On("FindAllArticle", 1, 10).Return(articles, nil).Once()
		repo.On("GetTotalArticleCount").Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetAllArticleUser(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_DeleteBookmarkArticle(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	userID := uint64(1)
	articleID := uint64(2)

	t.Run("Success Case - Bookmark Deleted", func(t *testing.T) {
		repo.On("IsArticleAlreadyBookmarked", userID, articleID).Return(true, nil).Once()
		repo.On("DeleteBookmarkArticle", userID, articleID).Return(nil).Once()

		err := service.DeleteBookmarkArticle(userID, articleID)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Article Not Bookmarked", func(t *testing.T) {
		repo.On("IsArticleAlreadyBookmarked", userID, articleID).Return(false, nil).Once()

		err := service.DeleteBookmarkArticle(userID, articleID)

		assert.Error(t, err)
		assert.Equal(t, "artikel tidak ditemukan", err.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Checking Bookmark", func(t *testing.T) {
		expectedErr := errors.New("gagal mengecek artikel")

		repo.On("IsArticleAlreadyBookmarked", userID, articleID).Return(false, expectedErr).Once()

		err := service.DeleteBookmarkArticle(userID, articleID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Deleting Bookmark", func(t *testing.T) {
		repo.On("IsArticleAlreadyBookmarked", userID, articleID).Return(true, nil).Once()

		expectedErr := errors.New("gagal menghapus artikel tersimpan")

		repo.On("DeleteBookmarkArticle", userID, articleID).Return(expectedErr).Once()

		err := service.DeleteBookmarkArticle(userID, articleID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestArticleService_GetUserBookmarkArticle(t *testing.T) {
	repo := mocks.NewRepositoryArticleInterface(t)
	service := NewArticleService(repo)

	userID := uint64(1)

	t.Run("Success Case - Success Get User Bookmark Articles", func(t *testing.T) {
		bookmarkArticles := []*entities.ArticleBookmarkModels{
			{ID: 1, UserID: 1, ArticleID: 1},
			{ID: 2, UserID: 2, ArticleID: 2},
		}

		repo.On("GetUserBookmarkArticle", userID).Return(bookmarkArticles, nil).Once()

		result, err := service.GetUserBookmarkArticle(userID)

		assert.NoError(t, err)
		assert.Equal(t, len(bookmarkArticles), len(result))
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Get User Bookmark Articles", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan artikel tersimpan user")

		repo.On("GetUserBookmarkArticle", userID).Return(nil, expectedErr).Once()

		result, err := service.GetUserBookmarkArticle(userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestGetFilterDateRange(t *testing.T) {
	service := &ArticleService{}

	tests := []struct {
		filterType    string
		expectedStart time.Time
		expectedEnd   time.Time
		expectedError error
	}{
		{
			filterType:    "minggu ini",
			expectedStart: time.Now().In(time.UTC).AddDate(0, 0, -int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour),
			expectedEnd:   time.Now().In(time.UTC).AddDate(0, 0, 7-int(time.Now().In(time.UTC).Weekday())).Truncate(24 * time.Hour),
			expectedError: nil,
		},
		{
			filterType:    "bulan ini",
			expectedStart: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
			expectedError: nil,
		},
		{
			filterType:    "tahun ini",
			expectedStart: time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(time.Now().Year()+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
			expectedError: nil,
		},
		{
			filterType:    "invalid",
			expectedError: errors.New("tipe filter tidak valid"),
		},
	}

	for _, test := range tests {
		startDate, endDate, err := service.GetFilterDateRange(test.filterType)

		if !startDate.Equal(test.expectedStart) || !endDate.Equal(test.expectedEnd) || !reflect.DeepEqual(err, test.expectedError) {
			t.Errorf("For filter type %s, expected (%v, %v, %v), but got (%v, %v, %v)", test.filterType, test.expectedStart, test.expectedEnd, test.expectedError, startDate, endDate, err)
		}
	}
}
