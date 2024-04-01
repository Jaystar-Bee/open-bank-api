package routes

import (
	"github.com/Jaystar-Bee/open-bank-api/handlers"
	"github.com/Jaystar-Bee/open-bank-api/middlewares"
	"github.com/gin-gonic/gin"
)

func WalletRoutes(server *gin.Engine) {

	walletRoutes := server.Group("/wallet").Use(middlewares.CheckAuthentication)
	walletRoutes.GET("", handlers.GetWallet)
	walletRoutes.POST("/send", handlers.SendToUser)
	walletRoutes.POST("/requests", handlers.RequestMoney)
	walletRoutes.GET("/requests", nil)
	walletRoutes.DELETE("/requests/:id", handlers.DeleteRequest)
	walletRoutes.POST("/requests/confirm", nil)
	walletRoutes.POST("/requests/reject", nil)
	walletRoutes.POST("/deposit", handlers.Deposit)
}
