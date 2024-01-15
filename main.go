package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	routes "github.com/shubham-2909/jwtAuth/router"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	routes.AuthRoutes(router)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello This is My first API Using gin")
	})

	router.Run(":" + port)

}
