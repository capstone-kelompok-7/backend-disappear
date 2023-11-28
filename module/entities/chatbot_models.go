package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatModel struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IdUser    string             `json:"id_user" form:"id_user"`
	Role      string
	Text      string `json:"text" form:"text"`
	CreatedAt time.Time
}
