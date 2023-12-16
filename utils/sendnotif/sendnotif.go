package sendnotif

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func SendNotification(title string, body string, token string) (string, error) {
	opt := option.WithCredentialsFile("contoh-c4760-firebase-adminsdk-i3rk9-0a642f465e.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Error("Error initializing Firebase app", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		logrus.Error("Error creating Firebase Messaging client", err)
	}

	message := &messaging.Message{
		Data: map[string]string{
			"order_id": request.OrderID,
			"title":    request.Title,
			"body":     request.Body,
		},
		Token: request.Token,
	}

	response, err := client.Send(context.Background(), message)
	if err != nil {
		if messaging.IsRegistrationTokenNotRegistered(err) {
			logrus.Error("Registration token is not valid. Remove it from your database.", err)
		} else {
			logrus.Error("Error sending message: ", err)
		}
	}

	logrus.Info("Successfully sent message:", response)
	return response, nil
}
