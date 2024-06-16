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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var courseCollection *mongo.Collection = database.OpenCollection(database.Client, "course")

func GetCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := courseCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing course items"})
			return
		}

		defer result.Close(ctx)

		var allCourses []bson.M
		if err = result.All(ctx, &allCourses); err != nil {
			log.Fatal(err)
		}

		if len(allCourses) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No courses found"})
			return
		}

		c.JSON(http.StatusOK, allCourses)
	}
}

func GetCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		courseId := c.Param("course_id")

		objID, errr := primitive.ObjectIDFromHex(courseId)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
			return
		}

		var course models.Course

		err := courseCollection.FindOne(ctx, bson.M{"_id": objID}, options.FindOne()).Decode(&course)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the courses"})
		}

		c.JSON(http.StatusOK, course)
	}
}

func CreateCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

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

		course.Created_at = time.Now()
		course.Updated_at = time.Now()
		course.ID = primitive.NewObjectID()

		result, insertErr := courseCollection.InsertOne(ctx, course)

		if insertErr != nil {
			msg := "Course item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		courseId := c.Param("course_id")

		objID, errr := primitive.ObjectIDFromHex(courseId)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
			return
		}

		var updateCourse models.Course
		if err := c.BindJSON(&updateCourse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var course models.Course
		err := courseCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&course)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "course not found"})
			return
		}

		if updateCourse.Academic_year != nil {
			course.Academic_year = updateCourse.Academic_year
		}
		if updateCourse.Class_code != nil {
			course.Class_code = updateCourse.Class_code
		}
		if updateCourse.Course_code != nil {
			course.Course_code = updateCourse.Course_code
		}
		if updateCourse.Degree != nil {
			course.Degree = updateCourse.Degree
		}
		if updateCourse.Department != nil {
			course.Department = updateCourse.Department
		}
		if updateCourse.Instructor != nil {
			course.Instructor = updateCourse.Instructor
		}
		if updateCourse.Name != nil {
			course.Name = updateCourse.Name
		}
		if updateCourse.Semester != nil {
			course.Semester = updateCourse.Semester
		}
		if updateCourse.English != nil {
			course.English = updateCourse.English
		}

		course.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		_, updateErr := courseCollection.ReplaceOne(ctx, bson.M{"_id": objID}, course)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while updating course"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "course updated successfully"})
	}
}

func DeleteCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		courseID := c.Param("course_id")

		objID, err := primitive.ObjectIDFromHex(courseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
			return
		}

		filter := bson.M{"_id": objID}

		result, err := courseCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while deleting the course item"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "course item not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "course item deleted successfully"})
	}
}
