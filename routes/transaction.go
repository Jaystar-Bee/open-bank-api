package routes

import (
	"github.com/Jaystar-Bee/open-bank-api/middlewares"
	"github.com/gin-gonic/gin"
)

func TransactionRoutes(server *gin.Engine) {
	server.GET("/transaction", middlewares.CheckAuthentication, nil)
}
