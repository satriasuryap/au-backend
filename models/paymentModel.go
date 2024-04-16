package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID           primitive.ObjectID `bson:"_id"`
	Payment_Date time.Time          `json:"payment_date" validate:"required"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	Payment_id   string             `json:"payment_id"`
	User_id      *string            `json:"user_id" validate:"required"`
}
