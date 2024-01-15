package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shubham-2909/jwtAuth/database"
	"github.com/shubham-2909/jwtAuth/helpers"
	"github.com/shubham-2909/jwtAuth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(enterdPassword string, actualPassword string) (isValid bool, msg string) {
	err := bcrypt.CompareHashAndPassword([]byte(actualPassword), []byte(enterdPassword))
	check := true
	message := ""
	if err != nil {
		check = false
		message = "Email or Password is Incorrect"
		return check, message
	}
	return check, message
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking for email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email  already exists"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		uid := user.ID.String()

		token, refreshToken := helpers.GenerateTokens(*user.Email, *user.Name, *user.UserType, uid)
		user.Token = &token
		user.Refresh_token = &refreshToken
		resultInsertionNum, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User was not Created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, resultInsertionNum)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Credentials!"})
			return
		}
		passwordInvalid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordInvalid {
			fmt.Println("Im going here")
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something Went Wrong"})
			return
		}
		uid := foundUser.ID.String()
		token, refreshToken := helpers.GenerateTokens(*foundUser.Email, *foundUser.Name, *foundUser.UserType, uid)

		helpers.UpdateAllTokens(token, refreshToken, foundUser.ID)
		err = userCollection.FindOne(ctx, bson.M{"_id": foundUser.ID}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}
