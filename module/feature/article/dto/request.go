package dto

type CreateArticleRequest struct {
	Title   string `form:"title" validate:"required"`
	Photo   string `form:"photo"`
	Content string `form:"content" validate:"required"`
}

type UpdateArticleRequest struct {
	Title   string `form:"title"`
	Photo   string `form:"photo"`
	Content string `form:"content"`
}

type UserBookmarkRequest struct {
	UserID    uint64 `form:"user_id" json:"user_id"`
	ArticleID uint64 `form:"article_id" json:"article_id" validate:"required"`
}
