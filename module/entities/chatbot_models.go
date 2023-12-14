package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatModel struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    uint64             `json:"user_id" form:"user_id"`
	Role      string             `json:"role" form:"role"`
	Text      string             `json:"text" form:"text"`
	CreatedAt time.Time          `json:"created_at" form:"created_at"`
}
