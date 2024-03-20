package routes

import "github.com/gin-gonic/gin"

func UserRoutes(server *gin.Engine) {
	server.POST("/signup")
}
