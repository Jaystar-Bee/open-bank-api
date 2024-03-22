package routes

import (
	"github.com/Jaystar-Bee/open-bank-api/handlers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(server *gin.Engine) {
	userRoute := server.Group("/user")
	userRoute.POST("/signup", handlers.CreateUser)
	userRoute.GET("/tag/:tag", handlers.GetUserByTag)
	userRoute.GET("/email/:email", handlers.GetUserByEmail)
	userRoute.GET("/phone/:phone", handlers.GetUserByPhone)
	userRoute.GET("/:id", handlers.GetUserById)
}
