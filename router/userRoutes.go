package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shubham-2909/jwtAuth/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/", controllers.GetUsers())
	incomingRoutes.GET("user/:id", controllers.GetUser())
}
