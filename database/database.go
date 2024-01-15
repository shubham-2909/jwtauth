package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	mongo_url := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongo_url))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

var Client *mongo.Client = GetDBInstance()

func OpenCollection(client *mongo.Client, colName string) *mongo.Collection {
	collection := client.Database("GolangJWT").Collection(colName)
	return collection
}
