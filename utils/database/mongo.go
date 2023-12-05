package database

import (
	"context"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitMongoDB(config config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(config.MongoURL)

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Fatal("Database : cannot connect to mongo atlas database ", err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.Fatal(err.Error())
		//logrus.Errorf("Error pinging MongoDB: %v", err)
	}

	logrus.Info("Connected to MongoDB Atlas")
	return client
}
