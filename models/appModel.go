package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Apps struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Name      *string            `json:"appname" gorm:"unique" bson:"username,omitempty"`
	Icon      *string            `json:"icon"`
	CreatedAt time.Time          `json:"createdat" bson:"createat"`
	UpdatedAt time.Time          `json:"updatedat" bson:"updatedat"`
	Apps_ID   string             `json:"apps_id"`
	Pref_id   *string            `json:"pref_id"`
}
