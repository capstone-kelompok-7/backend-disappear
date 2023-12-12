package repository

import (
	"context"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/chatbot"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatbotRepository struct {
	collection *mongo.Collection
}

func NewChatbotRepository(db *mongo.Client) chatbot.RepositoryChatbotInterface {
	collection := db.Database("chatbot").Collection("chats")

	return &ChatbotRepository{
		collection: collection,
	}
}

func (r *ChatbotRepository) CreateQuestion(newData entities.ChatModel) error {
	ctx := context.Background()
	if _, err := r.collection.InsertOne(ctx, newData); err != nil {
		return nil
	}
	return nil

}

func (r *ChatbotRepository) CreateAnswer(newData entities.ChatModel) error {
	ctx := context.Background()
	if _, err := r.collection.InsertOne(ctx, newData); err != nil {
		return nil
	}
	return nil
}

func (r *ChatbotRepository) GetChatByIdUser(id uint64) ([]entities.ChatModel, error) {
	res, err := r.collection.Find(context.Background(), bson.M{"userid": id})
	if err != nil {
		logrus.Error("cant get chat by id user", err.Error())
	}
	var chats []entities.ChatModel

	for res.Next(context.Background()) {
		var chat entities.ChatModel
		if err := res.Decode(&chat); err != nil {
			logrus.Error(err)
		}
		chats = append(chats, chat)
	}

	return chats, nil

}
