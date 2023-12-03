package dto

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
)

type ArticleFormatter struct {
	ID      uint64    `json:"id"`
	Title   string    `json:"title"`
	Photo   string    `json:"photo"`
	Content string    `json:"content"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	Views   uint64    `json:"views"`
}

func FormatArticle(article entities.ArticleModels) ArticleFormatter {
	articleFormatter := ArticleFormatter{}
	articleFormatter.ID = article.ID
	articleFormatter.Title = article.Title
	articleFormatter.Photo = article.Photo
	articleFormatter.Content = article.Content
	articleFormatter.Author = article.Author
	articleFormatter.Date = article.UpdatedAt
	articleFormatter.Views = article.Views

	return articleFormatter
}

func FormatterArticle(articles []entities.ArticleModels) []ArticleFormatter {
	var articleFormatter []ArticleFormatter

	for _, article := range articles {
		formatArticle := FormatArticle(article)
		articleFormatter = append(articleFormatter, formatArticle)
	}

	return articleFormatter
}

type BookmarkedArticleFormatter struct {
	Title   string    `json:"title"`
	Photo   string    `json:"photo"`
	Content string    `json:"content"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	Views   uint64    `json:"views"`
}

type UserBookmarksResponse struct {
	ID        uint64                     `json:"id"`
	UserID    uint64                     `json:"user_id"`
	ArticleID uint64                     `json:"voucher_id"`
	Article   BookmarkedArticleFormatter `json:"article"`
}

func UserBookmarkFormatter(userBookmark []*entities.UserBookmarkModels) ([]UserBookmarksResponse, error) {
	var userBookmarks []UserBookmarksResponse

	for _, bookmark := range userBookmark {
		userBookmark := UserBookmarksResponse{
			ID:        bookmark.ID,
			UserID:    bookmark.UserID,
			ArticleID: bookmark.ArticleID,
			Article: BookmarkedArticleFormatter{
				Title:   bookmark.Article.Title,
				Photo:   bookmark.Article.Photo,
				Content: bookmark.Article.Content,
				Author:  bookmark.Article.Author,
				Date:    bookmark.Article.UpdatedAt,
				Views:   bookmark.Article.Views,
			},
		}
		userBookmarks = append(userBookmarks, userBookmark)

	}

	return userBookmarks, nil
}
