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

var newsCollection *mongo.Collection = database.OpenCollection(database.Client, "news")

func GetNews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := newsCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing table items"})
		}
		var allNews []bson.M
		if err = result.All(ctx, &allNews); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allNews)
	}
}

func GetNew() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		newsId := c.Param("news_id")
		var news models.News

		err := newsCollection.FindOne(ctx, bson.M{"news_id": newsId}).Decode(&news)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the tables"})
		}
		c.JSON(http.StatusOK, news)
	}
}

func CreateNews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var news models.News

		if err := c.BindJSON(&news); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(news)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		/*pref.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		pref.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))*/

		news.ID = primitive.NewObjectID()
		news.News_id = news.ID.Hex()

		result, insertErr := newsCollection.InsertOne(ctx, news)

		if insertErr != nil {
			msg := fmt.Sprintf("news item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func UpdateNews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var news models.News

		newsId := c.Param("news_id")

		if err := c.BindJSON(&news); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if news.Content != nil {
			updateObj = append(updateObj, bson.E{Key: "content", Value: news.Content})
		}
		if news.Title != nil {
			updateObj = append(updateObj, bson.E{Key: "title", Value: news.Title})
		}

		//pref.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"news_id": newsId}

		result, err := newsCollection.UpdateOne(
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
