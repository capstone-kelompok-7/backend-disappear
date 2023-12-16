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

func getDecodedFireBaseKey() ([]byte, error) {
	fireBaseAuthKey := config.InitConfig().FirebaseKey
	decodedKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}

func SendNotification(request SendNotificationRequest) (string, error) {
	decodedKey, err := getDecodedFireBaseKey()

	if err != nil {
		return "", err
	}

	opt := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	app, err := firebase.NewApp(context.Background(), nil, opt...)

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
