package controllers

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shubham-2909/jwtAuth/database"
	"github.com/shubham-2909/jwtAuth/helpers"
	"github.com/shubham-2909/jwtAuth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
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

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			log.Panic(err)
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			log.Panic(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		l := int64(limit)
		p := int64(page)
		skip := (p*l - l)
		fOpt := options.FindOptions{Skip: &skip, Limit: &l}

		totalUsers, err := userCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Panic(err)
		}
		totalPages := math.Ceil(float64(totalUsers) / float64(limit))
		cursor, err := userCollection.Find(ctx, bson.M{}, &fOpt)
		if err != nil {
			log.Panic(err)
		}
		defer cursor.Close(ctx)
		var users []primitive.M
		for cursor.Next(ctx) {
			var user bson.M
			err := cursor.Decode(&user)
			if err != nil {
				log.Panic(err)
			}
			users = append(users, user)
		}
		c.JSON(http.StatusOK, gin.H{"users": users, "totalPages": totalPages})
	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ID := c.Param("id")
		if err := helpers.MapUserTypetoID(c, ID); err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}
		userId, _ := primitive.ObjectIDFromHex(ID)
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancel()
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
		if err != nil {
			log.Panic(err)
		}
		c.JSON(http.StatusOK, user)
	}
}
