package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transcript struct {
	ID            primitive.ObjectID `bson:"_id"`
	User_id       *string            `json:"user_id"`
	Course_id     *string            `json:"course_id"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	Transcript_id string             `json:"transcript_id"`
}
