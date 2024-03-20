package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prefs struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	User_id string             `json:"user_id"`
	Apps_id string             `json:"apps_id"`
}
