package routes

import (
	"github.com/Jaystar-Bee/open-bank-api/handlers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(server *gin.Engine) {
	userRoute := server.Group("/user")
	userRoute.POST("/signup", handlers.CreateUser)
}
