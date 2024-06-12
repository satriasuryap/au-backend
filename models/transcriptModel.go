package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transcript struct {
	ID         primitive.ObjectID `bson:"_id"`
	GPA        *float64           `json:"gpa"`
	Taken_in   *string            `json:"taken_in"`
	Course_id  *string            `json:"course_id"`
	User_id    *string            `json:"user_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
