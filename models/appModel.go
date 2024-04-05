package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Apps struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	Name    string             `json:"appname" gorm:"unique" bson:"username,omitempty"`
	Icon    *string            `json:"icon"`
	Apps_ID string             `json:"apps_id"`
}
