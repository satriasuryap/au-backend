package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type News struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Title         *string            `json:"title" gorm:"unique" bson:"title,omitempty"`
	Content       *string            `json:"content" gorm:"unique" bson:"content,omitempty"`
	URL           string             `json:"url"`
	CreatedAt     time.Time          `json:"createdat" bson:"createat"`
	DeactivatedAt time.Time          `json:"updatedat" bson:"updatedat"`
	News_id       string             `json:"news_id"`
}
