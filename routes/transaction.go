package routes

import (
	"github.com/Jaystar-Bee/open-bank-api/handlers"
	"github.com/Jaystar-Bee/open-bank-api/middlewares"
	"github.com/gin-gonic/gin"
)

func TransactionRoutes(server *gin.Engine) {
	server.GET("/transactions", middlewares.CheckAuthentication, handlers.GetTransactions)
	server.GET("/transactions/:id", middlewares.CheckAuthentication, handlers.GetTransactionByID)
}
