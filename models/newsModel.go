package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type News struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Title         string             `json:"username" gorm:"unique" bson:"username,omitempty"`
	Content       string             `json:"email" gorm:"unique" bson:"email,omitempty"`
	URL           string             `json:"url"`
	CreatedAt     time.Time          `json:"createdat" bson:"createat"`
	DeactivatedAt time.Time          `json:"updatedat" bson:"updatedat"`
}
