package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Username      string             `json:"username" gorm:"unique" bson:"username,omitempty"`
	Email         string             `json:"email" gorm:"unique" bson:"email,omitempty"`
	Password      []byte             `json:"password" bson:"password"`
	CreatedAt     time.Time          `json:"createdat" bson:"createat"`
	DeactivatedAt time.Time          `json:"updatedat" bson:"updatedat"`
}
