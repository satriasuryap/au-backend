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

var transcriptCollection *mongo.Collection = database.OpenCollection(database.Client, "transcript")

func GetTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		var transcripts []models.Transcript

		cursor, err := transcriptCollection.Find(ctx, bson.M{"user_id": userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the transcript"})
			return
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var transcript models.Transcript
			if err := cursor.Decode(&transcript); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding transcript"})
				return
			}
			transcripts = append(transcripts, transcript)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred during cursor iteration"})
			return
		}

		c.JSON(http.StatusOK, transcripts)
	}
}

func GetApprovedTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		var transcripts []models.Transcript

		filter := bson.M{"user_id": userId, "approval": true}

		cursor, err := transcriptCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the transcript"})
			return
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var transcript models.Transcript
			if err := cursor.Decode(&transcript); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding transcript"})
				return
			}
			transcripts = append(transcripts, transcript)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred during cursor iteration"})
			return
		}

		c.JSON(http.StatusOK, transcripts)
	}
}

func GetNotApprovedTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		var transcripts []models.Transcript

		filter := bson.M{"user_id": userId, "approval": false}

		cursor, err := transcriptCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the transcript"})
			return
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var transcript models.Transcript
			if err := cursor.Decode(&transcript); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding transcript"})
				return
			}
			transcripts = append(transcripts, transcript)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred during cursor iteration"})
			return
		}

		c.JSON(http.StatusOK, transcripts)
	}
}

func CreateTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		userObjID, errr := primitive.ObjectIDFromHex(userId)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		var user models.User
		errUser := userCollection.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
		if errUser != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not found"})
			return
		}

		courseId := c.Param("course_id")
		courseObjID, errr := primitive.ObjectIDFromHex(courseId)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
			return
		}

		var course models.Course
		errCourse := courseCollection.FindOne(ctx, bson.M{"_id": courseObjID}).Decode(&course)
		if errCourse != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "course was not found"})
			return
		}

		var transcript models.Transcript
		if errTranscript := c.BindJSON(&transcript); errTranscript != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errTranscript.Error()})
			return
		}

		var existingTranscript models.Transcript
		errExisting := transcriptCollection.FindOne(ctx, bson.M{
			"user_id":   userId,
			"course_id": courseId,
			"taken_in":  transcript.Taken_in,
		}).Decode(&existingTranscript)

		if errExisting == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user has already taken this course in the specified semester"})
			return
		}

		transcript.Course_id = &courseId
		transcript.User_id = &userId

		gpa := 0.0
		transcript.GPA = &gpa

		initialApproval := false
		transcript.Approval = &initialApproval

		transcript.Created_at = time.Now()
		transcript.Updated_at = time.Now()
		transcript.ID = primitive.NewObjectID()

		result, insertErr := transcriptCollection.InsertOne(ctx, transcript)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transcript item was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		transcriptId := c.Param("transcript_id")

		objID, errr := primitive.ObjectIDFromHex(transcriptId)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transcript ID"})
			return
		}

		var updateTranscript models.Transcript
		if err := c.BindJSON(&updateTranscript); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var transcript models.Transcript
		err := transcriptCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&transcript)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transcript not found"})
			return
		}

		if updateTranscript.GPA != nil {
			transcript.GPA = updateTranscript.GPA
		}
		if updateTranscript.Approval != nil {
			transcript.Approval = updateTranscript.Approval
		}

		transcript.Updated_at = time.Now()

		_, updateErr := transcriptCollection.ReplaceOne(ctx, bson.M{"_id": objID}, transcript)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while updating transcript"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "transcript updated successfully"})
	}
}

func UpdateApprovalTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var request struct {
			TranscriptIDs []string `json:"transcript_ids"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		objectIDs := make([]primitive.ObjectID, len(request.TranscriptIDs))
		for i, id := range request.TranscriptIDs {
			objectID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transcript ID format"})
				return
			}
			objectIDs[i] = objectID
		}

		filter := bson.M{"_id": bson.M{"$in": objectIDs}}
		update := bson.M{"$set": bson.M{"approval": true}}

		result, err := transcriptCollection.UpdateMany(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating transcripts"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transcripts updated successfully", "modified_count": result.ModifiedCount})
	}
}

func DeleteTranscript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		transcriptID := c.Param("transcript_id")

		objID, err := primitive.ObjectIDFromHex(transcriptID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transcript ID"})
			return
		}

		filter := bson.M{"_id": objID}

		result, err := transcriptCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while deleting the transcript item"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "transcript item not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "transcript item deleted successfully"})
	}
}
