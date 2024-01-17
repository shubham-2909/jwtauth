package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shubham-2909/jwtAuth/controllers"
	"github.com/shubham-2909/jwtAuth/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/", controllers.GetUsers())
	incomingRoutes.GET("user/:id", controllers.GetUser())
}
