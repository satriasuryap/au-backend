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

var AppsCollection *mongo.Collection = database.OpenCollection(database.Client, "apps")

func GetApps() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := AppsCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing table items"})
		}
		var allApps []bson.M
		if err = result.All(ctx, &allApps); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allApps)
	}
}

func GetApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		appsId := c.Param("apps_id")
		var apps models.Apps

		err := AppsCollection.FindOne(ctx, bson.M{"apps_id": appsId}).Decode(&apps)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the tables"})
		}
		c.JSON(http.StatusOK, apps)
	}
}

func CreateApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var apps models.Apps

		if err := c.BindJSON(&apps); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(apps)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		/*pref.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		pref.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))*/

		apps.ID = primitive.NewObjectID()
		apps.Apps_ID = apps.ID.Hex()

		result, insertErr := AppsCollection.InsertOne(ctx, apps)

		if insertErr != nil {
			msg := fmt.Sprintf("apps item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func UpdateApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var apps models.Apps

		appsId := c.Param("apps_id")

		if err := c.BindJSON(&apps); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if apps.Icon != nil {
			updateObj = append(updateObj, bson.E{Key: "icon", Value: apps.Icon})
		}

		//pref.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"apps_id": appsId}

		result, err := AppsCollection.UpdateOne(
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
