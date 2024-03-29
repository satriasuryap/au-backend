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

var prefsCollection *mongo.Collection = database.OpenCollection(database.Client, "prefs")

func GetPrefs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := userCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing table items"})
		}
		var allPrefs []bson.M
		if err = result.All(ctx, &allPrefs); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allPrefs)
	}
}

func GetPref() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		prefId := c.Param("pref_id")
		var pref models.Prefs

		err := prefsCollection.FindOne(ctx, bson.M{"pref_id": prefId}).Decode(&pref)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the tables"})
		}
		c.JSON(http.StatusOK, pref)
	}
}

func CreatePref() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var pref models.Prefs

		if err := c.BindJSON(&pref); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(pref)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		/*pref.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		pref.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))*/

		pref.ID = primitive.NewObjectID()
		pref.Pref_id = pref.ID.Hex()

		result, insertErr := prefsCollection.InsertOne(ctx, pref)

		if insertErr != nil {
			msg := fmt.Sprintf("Table item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func UpdatePref() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var pref models.Prefs

		prefId := c.Param("pref_id")

		if err := c.BindJSON(&pref); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if pref.Apps_id != nil {
			updateObj = append(updateObj, bson.E{Key: "applications", Value: pref.Apps_id})
		}

		//pref.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"pref_id": prefId}

		result, err := prefsCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("table item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
