package chatbot

import (
	"context"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

type RepositoryChatbotInterface interface {
	GetChatByIdUser(id string) ([]entities.ChatModel, error)
	CreateQuestion(chat entities.ChatModel) error
	CreateAnswer(chat entities.ChatModel) error
}

type ServicChatbotInterface interface {
	GetChatByIdUser(id string) ([]entities.ChatModel, error)
	CreateQuestion(chat entities.ChatModel) error
	CreateAnswer(chat entities.ChatModel, name string) (string, error)
	GetAnswerFromAi(chat []openai.ChatCompletionMessage, ctx context.Context) (openai.ChatCompletionResponse, error)
}

type HandlerChatbotInterface interface {
	GetChatByIdUser() echo.HandlerFunc
	CreateQuestion() echo.HandlerFunc
	CreateAnswer() echo.HandlerFunc
}
