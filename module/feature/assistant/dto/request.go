package dto

type GenerateArticleAiRequest struct {
	Text string `json:"text" form:"text"`
}

type CreateChatRequest struct {
	Text string `json:"text" form:"text"`
}
