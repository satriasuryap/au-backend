package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Username      string             `json:"username" gorm:"unique" bson:"username,omitempty"`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Email         string             `json:"email" gorm:"unique" bson:"email,omitempty"`
	Password      []byte             `json:"password" bson:"password"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_token"`
	CreatedAt     time.Time          `json:"createdat" bson:"createat"`
	DeactivatedAt time.Time          `json:"updatedat" bson:"updatedat"`
	User_id       string             `json:"user_id"`
}
