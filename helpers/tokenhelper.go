package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/shubham-2909/jwtAuth/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	Name      string
	User_type string
	Uid       string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GenerateTokens(email string, userName string, user_type string, uid string) (signedToken string, signedRefreshToken string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}
	claims := &SignedDetails{
		Email:     email,
		User_type: user_type,
		Name:      userName,
		Uid:       uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 167)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	tokenSecret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	fmt.Println(tokenSecret)
	fmt.Println(refreshSecret)
	tokentoReturn, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Panic(err)
	}
	refreshToken, err := refresh_token.SignedString([]byte(refreshSecret))
	if err != nil {
		log.Panic(err)
	}
	return tokentoReturn, refreshToken
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId primitive.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var upDateObj []primitive.M
	upDateObj = append(upDateObj, bson.M{"$set": bson.M{"token": signedToken}})
	upDateObj = append(upDateObj, bson.M{"$set": bson.M{"refresh_token": signedRefreshToken}})
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	upDateObj = append(upDateObj, bson.M{"$set": bson.M{"updated_at": Updated_at}})

	upsert := true
	filter := bson.M{"_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(ctx, filter, upDateObj, &opt)
	if err != nil {
		fmt.Println("Error in update token")
		log.Panic(err)
		return
	}
	fmt.Println("Tokens side is finished")
}
