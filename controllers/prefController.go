package controllers

import (
	"context"
	"golang-au-backend/database"
	"golang-au-backend/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var prefsCollection *mongo.Collection = database.OpenCollection(database.Client, "prefs")

func GetPrefs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := prefsCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the menu items"})
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
		prefId := c.Param("prefs_id")
		var prefs models.Prefs

		err := appCollection.FindOne(ctx, bson.M{"pref_id": prefId}).Decode(&prefs)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the menu"})
		}
		c.JSON(http.StatusOK, prefs)
	}
}

func CreatePref() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pref models.Prefs
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&pref); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(pref)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		pref.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		pref.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		pref.ID = primitive.NewObjectID()
		pref.Pref_id = pref.ID.Hex()

		result, insertErr := prefsCollection.InsertOne(ctx, pref)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu item was not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()
	}
}
