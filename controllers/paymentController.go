package controllers

import (
	"context"
	"fmt"
	"golang-au-backend/database"
	"golang-au-backend/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

var paymentCollection *mongo.Collection = database.OpenCollection(database.Client, "payment")

func GetPayments() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := paymentCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})
		}
		var allPayment []bson.M
		if err = result.All(ctx, &allPayment); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allPayment)
	}
}

func GetPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		paymentId := c.Param("payment_id")
		var payment models.Payment

		err := paymentCollection.FindOne(ctx, bson.M{"payment_id": paymentId}).Decode(&payment)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the payment"})
		}
		c.JSON(http.StatusOK, payment)
	}
}

func CreatePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var payment models.Payment

		if err := c.BindJSON(&payment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(payment)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if payment.User_id != nil {
			err := userCollection.FindOne(ctx, bson.M{"user_id": payment.User_id}).Decode(&user)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message:User was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
		}

		payment.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		payment.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		payment.ID = primitive.NewObjectID()
		payment.Payment_id = payment.ID.Hex()

		result, insertErr := paymentCollection.InsertOne(ctx, payment)

		if insertErr != nil {
			msg := fmt.Sprintf("order item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdatePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payment models.Payment

		var updateObj primitive.D

		paymentId := c.Param("payment_id")
		if err := c.BindJSON(&payment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payment.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", payment.Updated_at})

		upsert := true

		filter := bson.M{"payment_id": paymentId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := paymentCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$st", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("order item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
