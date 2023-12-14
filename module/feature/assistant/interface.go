package assistant

import (
	"context"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

type RepositoryAssistantInterface interface {
	GetChatByIdUser(id uint64) ([]entities.ChatModel, error)
	CreateQuestion(chat entities.ChatModel) error
	CreateAnswer(chat entities.ChatModel) error
	GetLastOrdersByUserID(userID uint64) ([]*entities.OrderModels, error)
	GetTopSellingProducts() ([]string, error)
	GetTopRatedProducts() ([]*entities.ProductModels, error)
}

type ServiceAssistantInterface interface {
	GetChatByIdUser(id uint64) ([]entities.ChatModel, error)
	CreateQuestion(userID uint64, newData entities.ChatModel) error
	CreateAnswer(userID uint64, newData entities.ChatModel) (string, error)
	GetAnswerFromAi(chat []openai.ChatCompletionMessage, ctx context.Context) (openai.ChatCompletionResponse, error)
	GenerateArticle(title string) (string, error)
	GenerateRecommendationProduct(userID uint64) ([]string, error)
}

type HandlerAssistantInterface interface {
	GetChatByIdUser() echo.HandlerFunc
	CreateQuestion() echo.HandlerFunc
	CreateAnswer() echo.HandlerFunc
	GenerateArticle() echo.HandlerFunc
	GetProductByIdUser() echo.HandlerFunc
}
