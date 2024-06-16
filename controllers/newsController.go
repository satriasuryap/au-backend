package controllers

import (
	"context"
	"golang-au-backend/database"
	"golang-au-backend/models"
	"log"
	"net/http"
	"time"

	"strconv"

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
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		skip := (page - 1) * recordPerPage

		result, err := newsCollection.Find(ctx, bson.D{}, options.Find().SetLimit(int64(recordPerPage)).SetSkip(int64(skip)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing news items"})
			return
		}

		var allNews []bson.M
		if err = result.All(ctx, &allNews); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing news items"})
			return
		}

		c.JSON(http.StatusOK, allNews)
	}
}

func GetNewsOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		newsId := c.Param("news_id")

		objID, errr := primitive.ObjectIDFromHex(newsId)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		var news models.News

		err := newsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&news)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusOK, nil)
			} else {
				log.Printf("Error occurred while fetching the payment: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching news"})
			}
			return
		}

		c.JSON(http.StatusOK, news)
	}
}

func CreateNews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var news models.News

		if err := c.BindJSON(&news); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(news); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		news.ID = primitive.NewObjectID()
		news.CreatedAt = time.Now()
		news.DeactivatedAt = time.Time{}

		resultInsertionNumber, insertErr := newsCollection.InsertOne(ctx, news)

		if insertErr != nil {
			msg := "news item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func UpdateNews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		newsId := c.Param("news_id")

		var updateNews models.News
		if err := c.BindJSON(&updateNews); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var news models.News
		err := newsCollection.FindOne(ctx, bson.M{"newsid": newsId}).Decode(&news)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "news not found"})
			return
		}

		if updateNews.Title != nil {
			news.Title = updateNews.Title
		}
		if updateNews.Content != nil {
			news.Content = updateNews.Content
		}
		if updateNews.NavigateTo != nil {
			news.NavigateTo = updateNews.NavigateTo
		}

		_, updateErr := newsCollection.ReplaceOne(ctx, bson.M{"newsid": newsId}, news)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while updating news"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "news updated successfully"})
	}
}

func DeleteNews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		newsId := c.Param("news_id")

		objID, err := primitive.ObjectIDFromHex(newsId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid news ID"})
			return
		}

		filter := bson.M{"_id": objID}

		result, err := newsCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while deleting the news item"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "news item not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "news item deleted successfully"})
	}
}
