package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	MongoDB := "mongodb://localhost:27017"
	fmt.Print(MongoDB)

	client, err := mongo.NewClient(options.NewClient)
}
