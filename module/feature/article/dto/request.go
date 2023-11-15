package dto

type CreateArticleRequest struct {
	Title     string `form:"title"`
	Photo     string `form:"photo"`
	Content   string `form:"content"`
}