package service

import (
	"errors"
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/mocks"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestAssistantService_CreateQuestion(t *testing.T) {
	repo := new(mocks.RepositoryAssistantInterface)
	openaiClient := new(openai.Client)
	config := config.Config{}

	service := NewAssistantService(repo, openaiClient, config)
	userID := uint64(1)
	question := entities.ChatModel{
		UserID:    userID,
		Role:      "question",
		Text:      "What is your question?",
		CreatedAt: time.Now(),
	}

	t.Run("Success Case - Create Question", func(t *testing.T) {
		repo.On("CreateQuestion", question).Return(nil).Once()

		err := service.CreateQuestion(userID, question)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

}

func TestAssistantService_GetChatByIdUser(t *testing.T) {
	repo := new(mocks.RepositoryAssistantInterface)
	openaiClient := new(openai.Client)
	config := config.Config{}

	service := NewAssistantService(repo, openaiClient, config)

	userID := uint64(1)
	id, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%x", userID))

	expectedChats := []entities.ChatModel{
		{ID: id, UserID: userID, Text: "Hi there!"},
		{ID: id, UserID: userID, Text: "How can I help you?"},
	}
	t.Run("Failed Case - Invalid UserID", func(t *testing.T) {
		var expectedErr = errors.New("failed to fetch chat by ID")
		repo.On("GetChatByIdUser", userID).Return(nil, expectedErr).Once()

		chats, err := service.GetChatByIdUser(userID)

		assert.Nil(t, chats)
		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
		repo.AssertExpectations(t)
	})

	t.Run("Success Case - Get Chat", func(t *testing.T) {
		repo.On("GetChatByIdUser", userID).Return(expectedChats, nil)

		chats, err := service.GetChatByIdUser(userID)

		assert.NoError(t, err)
		assert.NotNil(t, chats)
		assert.Equal(t, len(expectedChats), len(chats))
		for i := range expectedChats {
			assert.Equal(t, expectedChats[i].ID, chats[i].ID)
			assert.Equal(t, expectedChats[i].UserID, chats[i].UserID)
			assert.Equal(t, expectedChats[i].Text, chats[i].Text)
		}

		repo.AssertExpectations(t)
	})

}
