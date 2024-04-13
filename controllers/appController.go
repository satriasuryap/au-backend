package controllers

import (
	"context"
	"fmt"
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
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"app_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
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
		var prefs models.Prefs
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
		err := prefsCollection.FindOne(ctx, bson.M{"menu_id": apps.Pref_id}).Decode(&prefs)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		apps.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		apps.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		apps.ID = primitive.NewObjectID()
		apps.Apps_ID = apps.ID.Hex()

		result, insertErr := appCollection.InsertOne(ctx, apps)
		if insertErr != nil {
			msg := fmt.Sprintf("Apps item was not created")
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
		var prefs models.Prefs
		var apps models.Apps

		appsId := c.Param("apps_id")

		if err := c.BindJSON(&apps); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if apps.Name != nil {
			updateObj = append(updateObj, bson.E{"name", apps.Name})
		}

		if apps.Icon != nil {
			updateObj = append(updateObj, bson.E{"apps_icon", apps.Icon})
		}

		if apps.Pref_id != nil {
			err := prefsCollection.FindOne(ctx, bson.M{"pref_id": apps.Pref_id}).Decode(&prefs)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message:Menu was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{"prefs", apps.Name})
		}

		apps.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updatedat", apps.UpdatedAt})

		upsert := true
		filter := bson.M{"apps_id": appsId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := appCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("apps item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
