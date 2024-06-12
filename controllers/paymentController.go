package controllers

import (
	"context"
	"golang-au-backend/database"
	"golang-au-backend/models"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var paymentCollection *mongo.Collection = database.OpenCollection(database.Client, "payment")

func GetPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		var payment models.Payment

		err := paymentCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&payment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the payment"})
		}

		c.JSON(http.StatusOK, payment)
	}
}

func CreatePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		objID, errr := primitive.ObjectIDFromHex(userId)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		var payment models.Payment
		if err := c.BindJSON(&payment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payment.User_id = &userId

		validationErr := validate.Struct(payment)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not found"})
			return
		}

		var existingPayment models.Payment
		err = paymentCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&existingPayment)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "payment for this user already exists"})
			return
		} else if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking for existing payment"})
			return
		}

		payment.Created_at = time.Now()
		payment.Updated_at = time.Now()
		payment.ID = primitive.NewObjectID()

		result, insertErr := paymentCollection.InsertOne(ctx, payment)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "payment item was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdatePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		paymentId := c.Param("payment_id")
		objID, errr := primitive.ObjectIDFromHex(paymentId)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		}

		var updatePayment models.Payment
		if err := c.BindJSON(&updatePayment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var payment models.Payment
		err := paymentCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&payment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "payment not found"})
			return
		}

		if !updatePayment.Payment_Date.IsZero() {
			payment.Payment_Date = updatePayment.Payment_Date
		}
		if updatePayment.Primary_tuition != nil {
			payment.Primary_tuition = updatePayment.Primary_tuition
		}
		if updatePayment.Additional_tuition != nil {
			payment.Additional_tuition = updatePayment.Additional_tuition
		}
		if updatePayment.Scholarship != nil {
			payment.Scholarship = updatePayment.Scholarship
		}
		if updatePayment.Status != nil {
			payment.Status = updatePayment.Status
		}

		payment.Updated_at = time.Now()

		_, updateErr := paymentCollection.ReplaceOne(ctx, bson.M{"_id": objID}, payment)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while updating payment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "payment updated successfully"})
	}
}

func DeletePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		paymentID := c.Param("payment_id")

		objID, err := primitive.ObjectIDFromHex(paymentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
			return
		}

		filter := bson.M{"_id": objID}

		result, err := paymentCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while deleting the payment item"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment item not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "payment item deleted successfully"})
	}
}
