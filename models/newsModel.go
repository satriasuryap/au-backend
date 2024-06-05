package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type News struct {
	ID            primitive.ObjectID `bson:"_id"`
	Title         *string            `json:"title"`
	Content       *string            `json:"content"`
	NavigateTo    string             `json:"navigateto"`
	CreatedAt     time.Time          `json:"createdat"`
	DeactivatedAt time.Time          `json:"deactivatedat"`
	Newsid        string             `json:"newsid"`
}
