package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Payment_Date       time.Time          `json:"payment_date"`
	Created_at         time.Time          `json:"created_at"`
	Updated_at         time.Time          `json:"updated_at"`
	Primary_tuition    *float64           `json:"primary_tuition"`
	Additional_tuition *float64           `json:"additional_tuition"`
	Scholarship        *float64           `json:"scholarship"`
	Status             *string            `json:"status"`
	User_id            *string            `json:"user_id" validate:"required"`
}
