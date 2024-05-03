package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct{
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Name		  *string 			 `json:"name" bson:"name"`
	Class		  *string 			 `json:"class" bson:"class"`
	Instructor	  *string 			 `json:"instructor" bson:"instructor"`
	Course_id	  string 			 `json:"course_id"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}