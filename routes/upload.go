package routes

import (
	"github.com/Jaystar-Bee/open-bank-api/handlers"
	"github.com/Jaystar-Bee/open-bank-api/middlewares"
	"github.com/gin-gonic/gin"
)

func UploadRoutes(server *gin.Engine) {
	server.POST("/upload", middlewares.CheckAuthentication, handlers.UploadFile)
}
