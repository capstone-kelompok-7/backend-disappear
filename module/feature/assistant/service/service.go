package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type AssistantService struct {
	repo   assistant.RepositoryAssistantInterface
	openai *openai.Client
	debug  bool
	config config.Config
}

func NewAssistantService(repo assistant.RepositoryAssistantInterface, openai *openai.Client, config config.Config) assistant.ServiceAssistantInterface {
	return &AssistantService{
		repo:   repo,
		openai: openai,
		config: config,
		debug:  false,
	}
}

func (s *AssistantService) CreateQuestion(userID uint64, newData entities.ChatModel) error {
	value := &entities.ChatModel{
		UserID:    userID,
		Role:      "question",
		Text:      newData.Text,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateQuestion(*value); err != nil {
		return err
	}
	return nil
}

func (s *AssistantService) GetAnswerFromAi(chat []openai.ChatCompletionMessage, ctx context.Context) (openai.ChatCompletionResponse, error) {
	model := openai.GPT3Dot5Turbo
	resp, err := s.openai.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: chat,
		},
	)

	return resp, err
}

func (s *AssistantService) CreateAnswer(userID uint64, newData entities.ChatModel) (string, error) {
	ctx := context.Background()
	chat := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Kamu Adalah Chatbot Bertema Lingkungan",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Hi Bisa Bantu Saya Menjawab Pertanyaan Tentang Lingkungan",
		},
	}

	if newData.Text != "" {
		chat = append(chat, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Kamu akan di berikan sebuah pertanyaan mengenai %s, berikan jawabannya maksimal 20 kata", newData.Text),
		})
	}

	resp, err := s.GetAnswerFromAi(chat, ctx)
	if err != nil {
		logrus.Error("Can't Get Answer From Ai: ", err.Error())
	}

	if s.debug {
		fmt.Printf(
			"ID: %s. Created: %d. Model: %s. Choices: %v.\n",
			resp.ID, resp.Created, resp.Model, resp.Choices,
		)
	}

	answer := openai.ChatCompletionMessage{
		Role:    resp.Choices[0].Message.Role,
		Content: resp.Choices[0].Message.Content,
	}

	answerText := fmt.Sprintf(resp.Choices[0].Message.Content)

	value := &entities.ChatModel{
		UserID:    userID,
		Role:      "answer",
		Text:      answer.Content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateAnswer(*value); err != nil {
		logrus.Error("Can't create answer in the repository: ", err.Error())
		return "", err
	}
	return answerText, nil
}

func (s *AssistantService) GetChatByIdUser(id uint64) ([]entities.ChatModel, error) {
	res, err := s.repo.GetChatByIdUser(id)
	if err != nil {
		logrus.Error("error getting chat by id", err.Error())
	}
	return res, nil
}

func (s *AssistantService) GenerateArticle(title string) (string, error) {
	ctx := context.Background()
	chat := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Kamu Adalah Chatbot Bertema Lingkungan",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Hi Bisa Bantu Saya Menjawab Pertanyaan Tentang Lingkungan",
		},
	}

	if title != "" {
		chat = append(chat, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Buatlah Artikel Dengan Judul %s", title),
		})
	}

	resp, err := s.GetAnswerFromAi(chat, ctx)
	if err != nil {
		logrus.Error("Can't Get Answer From Ai: ", err.Error())
	}

	if s.debug {
		fmt.Printf(
			"ID: %s. Created: %d. Model: %s. Choices: %v.\n",
			resp.ID, resp.Created, resp.Model, resp.Choices,
		)
	}

	answer := openai.ChatCompletionMessage{
		Role:    resp.Choices[0].Message.Role,
		Content: resp.Choices[0].Message.Content,
	}

	return answer.Content, nil
}

func (s *AssistantService) GenerateRecommendationProduct(userID uint64) ([]string, error) {
	ctx := context.Background()

	orders, err := s.repo.GetLastOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	var recommendedProducts []string

	if len(orders) == 0 {
		topProducts, err := s.repo.GetTopRatedProducts()
		if err != nil {
			return nil, err
		}

		for _, product := range topProducts {
			recommendedProducts = append(recommendedProducts, product.Name)
		}
	} else {
		chat := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "You are an analyst for the user's purchases.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Based on the user's previous purchases, provide 3 relevant products. Just the product names, no need for description or others.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "example answers\n1. Tas\n2. Alat Makan\n3. Korek Api",
			},
		}

		orderContent := "List of products purchased by the user: \n"
		for _, order := range orders {
			for _, product := range order.OrderDetails {
				orderContent += fmt.Sprintf("- %s\n", product.Product.Name)
			}
		}

		chat = append(chat, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: orderContent,
		})

		resp, err := s.GetAnswerFromAi(chat, ctx)
		if err != nil {
			return nil, err
		}
		for _, choice := range resp.Choices {
			if choice.Message.Role == "assistant" {
				lines := strings.Split(choice.Message.Content, "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						product := strings.SplitN(line, ". ", 2)
						if len(product) > 1 {
							recommendedProducts = append(recommendedProducts, strings.TrimSpace(product[1]))
						}
					}
				}
			}
		}
	}

	return recommendedProducts, nil
}
