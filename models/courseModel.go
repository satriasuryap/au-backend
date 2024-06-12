package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Academic_year *string            `json:"academic_year" bson:"academic_year"`
	Degree        *string            `json:"degree" bson:"degree"`
	Department    *string            `json:"department" bson:"department"`
	Name          *string            `json:"name" bson:"name"`
	Class_code    *string            `json:"class_code" bson:"class_code"`
	Instructor    *string            `json:"instructor" bson:"instructor"`
	Course_code   *string            `json:"course_code" bson:"course_code"`
	Semester      *string            `json:"semester" bson:"semester"`
	English       *bool              `json:"english" bson:"english"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}
