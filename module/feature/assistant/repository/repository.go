package repository

import (
	"context"
	"gorm.io/gorm"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AssistantRepository struct {
	collection *mongo.Collection
	dbo        *gorm.DB
}

func NewAssistantRepository(db *mongo.Client, dbo *gorm.DB) assistant.RepositoryAssistantInterface {
	collection := db.Database("assistant").Collection("chats")

	return &AssistantRepository{
		collection: collection,
		dbo:        dbo,
	}
}

func (r *AssistantRepository) CreateQuestion(newData entities.ChatModel) error {
	ctx := context.Background()
	if _, err := r.collection.InsertOne(ctx, newData); err != nil {
		return nil
	}
	return nil

}

func (r *AssistantRepository) CreateAnswer(newData entities.ChatModel) error {
	ctx := context.Background()
	if _, err := r.collection.InsertOne(ctx, newData); err != nil {
		return nil
	}
	return nil
}

func (r *AssistantRepository) GetChatByIdUser(id uint64) ([]entities.ChatModel, error) {
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

func (r *AssistantRepository) GetLastOrdersByUserID(userID uint64) ([]*entities.OrderModels, error) {
	var orders []*entities.OrderModels

	if err := r.dbo.
		Preload("OrderDetails").
		Preload("OrderDetails.Product").
		Preload("OrderDetails.Product.ProductPhotos").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Limit(10).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *AssistantRepository) GetTopSellingProducts() ([]string, error) {
	var topProducts []string

	if err := r.dbo.
		Table("order_details").
		Select("product_name").
		Group("product_name").
		Order("COUNT(product_name) DESC").
		Limit(3).
		Pluck("product_name", &topProducts).
		Error; err != nil {
		return nil, err
	}

	return topProducts, nil
}

func (r *AssistantRepository) GetTopRatedProducts() ([]*entities.ProductModels, error) {
	var products []*entities.ProductModels

	if err := r.dbo.Order("rating desc").Where("deleted_at IS NULL").Limit(3).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
