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

func FormatArticle(article *entities.ArticleModels) *ArticleFormatter {
	articleFormatter := &ArticleFormatter{}
	articleFormatter.ID = article.ID
	articleFormatter.Title = article.Title
	articleFormatter.Photo = article.Photo
	articleFormatter.Content = article.Content
	articleFormatter.Author = article.Author
	articleFormatter.Date = article.CreatedAt
	articleFormatter.Views = article.Views

	return articleFormatter
}

func FormatterArticle(articles []*entities.ArticleModels) []*ArticleFormatter {
	var articleFormatter []*ArticleFormatter

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
	ArticleID uint64                     `json:"article_id"`
	Article   BookmarkedArticleFormatter `json:"article"`
}

func UserBookmarkFormatter(userBookmarks []*entities.ArticleBookmarkModels) ([]UserBookmarksResponse, error) {
	var formattedUserBookmarks []UserBookmarksResponse

	for _, bookmark := range userBookmarks {
		formattedBookmark := UserBookmarksResponse{
			ID:        bookmark.ID,
			UserID:    bookmark.UserID,
			ArticleID: bookmark.ArticleID,
			Article: BookmarkedArticleFormatter{
				Title:   bookmark.Article.Title,
				Photo:   bookmark.Article.Photo,
				Content: bookmark.Article.Content,
				Author:  bookmark.Article.Author,
				Date:    bookmark.Article.CreatedAt,
				Views:   bookmark.Article.Views,
			},
		}
		formattedUserBookmarks = append(formattedUserBookmarks, formattedBookmark)
	}

	return formattedUserBookmarks, nil
}
