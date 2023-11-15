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
}

func FormatArticle(article entities.ArticleModels) ArticleFormatter {
	articleFormatter := ArticleFormatter{}
	articleFormatter.ID = article.ID
	articleFormatter.Title = article.Title
	articleFormatter.Photo = article.Photo
	articleFormatter.Content = article.Content
	articleFormatter.Author = article.Author
	articleFormatter.Date = article.UpdatedAt

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
