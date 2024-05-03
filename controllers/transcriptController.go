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
)

var transcriptCollection *mongo.Collection = database.OpenCollection(database.Client, "transcript")

func GetTranscripts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := prefsCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the transcript items"})
		}
		var allTranscript []bson.M
		if err = result.All(ctx, &allTranscript); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allTranscript)
	}
}

func GetTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		transcriptId := c.Param("transcript_id")
		var transcript models.Transcript

		err := appCollection.FindOne(ctx, bson.M{"transcript_id": transcriptId}).Decode(&transcript)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the transcript"})
		}
		c.JSON(http.StatusOK, transcript)
	}
}

func CreateTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var transcript models.Transcript
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.BindJSON(&transcript); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(transcript)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		transcript.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		transcript.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		transcript.ID = primitive.NewObjectID()
		transcript.Transcript_id = transcript.ID.Hex()

		result, insertErr := prefsCollection.InsertOne(ctx, transcript)
		if insertErr != nil {
			msg := fmt.Sprintf("Transcript item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()
	}

}
