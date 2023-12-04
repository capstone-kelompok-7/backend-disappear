package service

import (
	"context"
	"fmt"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/chatbot"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type ChatbotService struct {
	repo   chatbot.RepositoryChatbotInterface
	openai *openai.Client
	debug  bool
	config config.Config
}

func NewChatbotService(repo chatbot.RepositoryChatbotInterface, openai *openai.Client, config config.Config) chatbot.ServicChatbotInterface {
	return &ChatbotService{
		repo:   repo,
		openai: openai,
		config: config,
		debug:  false,
	}
}

func (s *ChatbotService) CreateQuestion(newData entities.ChatModel) error {
	value := &entities.ChatModel{
		IdUser:    newData.IdUser,
		Role:      "question",
		Text:      newData.Text,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateQuestion(*value); err != nil {
		return err
	}
	return nil
}

func (s *ChatbotService) GetAnswerFromAi(chat []openai.ChatCompletionMessage, ctx context.Context) (openai.ChatCompletionResponse, error) {
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

func (s *ChatbotService) CreateAnswer(newData entities.ChatModel) (string, error) {
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
			Content: newData.Text,
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

	value := &entities.ChatModel{
		IdUser:    newData.IdUser,
		Role:      "answer",
		Text:      answer.Content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateAnswer(*value); err != nil {
		logrus.Error("Can't create answer in the repository: ", err.Error())
		return "", err
	}
	return answer.Content, nil
}

func (s *ChatbotService) GetChatByIdUser(id string) ([]entities.ChatModel, error) {
	res, err := s.repo.GetChatByIdUser(id)
	if err != nil {
		logrus.Error("error getting chat by id", err.Error())
	}
	return res, nil
}

func (s *ChatbotService) GenerateArtikelAi(judul string) (string, error) {
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

	if judul != "" {
		chat = append(chat, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Buatlah Artikel Dengan Judul %s", judul),
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
