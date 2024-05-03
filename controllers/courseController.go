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

var courseCollection *mongo.Collection = database.OpenCollection(database.Client, "course")

func GetCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := transcriptCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing course items"})
		}
		var allCourses []bson.M
		if err = result.All(ctx, &allCourses); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allCourses)
	}
}

func GetCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		courseId := c.Param("course_id")
		var course models.Course

		err := courseCollection.FindOne(ctx, bson.M{"course_id": courseId}).Decode(&course)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the courses"})
		}
		c.JSON(http.StatusOK, course)
	}
}

func CreateCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var course models.Course

		if err := c.BindJSON(&course); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(course)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		course.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		course.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		course.ID = primitive.NewObjectID()
		course.Course_id = course.ID.Hex()

		result, insertErr := courseCollection.InsertOne(ctx, course)

		if insertErr != nil {
			msg := fmt.Sprintf("Course item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}

func UpdateCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var course models.Course

		courseId := c.Param("course_id")

		if err := c.BindJSON(&course); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if course.Class != nil {
			updateObj = append(updateObj, bson.E{"class", course.Class})
		}

		if course.Name != nil {
			updateObj = append(updateObj, bson.E{"name", course.Name})
		}

		if course.Instructor != nil {
			updateObj = append(updateObj, bson.E{"instructor", course.Instructor})
		}

		course.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"course_id": courseId}

		result, err := courseCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("course item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
