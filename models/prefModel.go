package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prefs struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	User_id    *string            `json:"user_id"`
	Start_Date *time.Time         `json:"start_date"`
	End_Date   *time.Time         `json:"end_date"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Apps_id    *string            `json:"apps_id"`
	Pref_id    string             `json:"pref_id"`
}
