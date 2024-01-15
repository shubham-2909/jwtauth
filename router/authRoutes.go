package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shubham-2909/jwtAuth/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controllers.SignUp())
	incomingRoutes.POST("/login", controllers.Login())
}
