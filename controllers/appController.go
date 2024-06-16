package controllers

import (
	"context"
	"golang-au-backend/database"
	"golang-au-backend/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var appCollection *mongo.Collection = database.OpenCollection(database.Client, "app")
var validate = validator.New()

func GetApps() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, _ = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1},
					{Key: "app_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

		result, err := appCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing food items"})
		}
		var allApps []bson.M
		if err = result.All(ctx, &allApps); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allApps[0])
	}
}

func GetApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("apps_id")
		var food models.Apps

		err := appCollection.FindOne(ctx, bson.M{"apps_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
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

		apps.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		apps.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		apps.ID = primitive.NewObjectID()
		apps.Apps_ID = apps.ID.Hex()

		result, insertErr := appCollection.InsertOne(ctx, apps)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Apps item was not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var apps models.Apps

		appsId := c.Param("apps_id")

		if err := c.BindJSON(&apps); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if apps.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: apps.Name})
		}

		if apps.Icon != nil {
			updateObj = append(updateObj, bson.E{Key: "apps_icon", Value: apps.Icon})
		}

		apps.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updatedat", Value: apps.UpdatedAt})

		upsert := true
		filter := bson.M{"apps_id": appsId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := appCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "apps item update failed"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
