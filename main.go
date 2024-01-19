package main

import (
	"os"

	"github.com/gin-gonic/gin"
	routes "github.com/shubham-2909/jwtAuth/router"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	router.Run(":" + port)
}
