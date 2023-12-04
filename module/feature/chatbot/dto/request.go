package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRequest struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IdUser    string             `json:"id_user" form:"id_user"`
	Role      string
	Name      string `json:"name" form:"name"`
	Text      string `json:"text" form:"text"`
	CreatedAt time.Time
}
type GenerateArtikelAiRequest struct {
	Text string `json:"text" form:"text"`
}
