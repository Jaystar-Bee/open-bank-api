package routes

import (
	"github.com/Jaystar-Bee/open-bank-api/handlers"
	"github.com/Jaystar-Bee/open-bank-api/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(server *gin.Engine) {
	userRoute := server.Group("/user")
	userRoute.POST("/signup", handlers.CreateUser)
	userRoute.GET("/tag/:tag", handlers.GetUserByTag)
	userRoute.GET("/email/:email", handlers.GetUserByEmail)
	userRoute.GET("/phone/:phone", handlers.GetUserByPhone)
	userRoute.GET("/:id", handlers.GetUserById)
	userRoute.POST("/login", handlers.Login)
	userRoute.POST("/verify", handlers.VerifyAccount)
	userRoute.POST("/sendotp", handlers.SendOTP)
	userRoute.GET("/renew", middlewares.CheckAuthentication, handlers.RenewToken)
	userRoute.PUT("/edit", middlewares.CheckAuthentication, middlewares.CheckAccountActivation, handlers.EditUser)
	// userRoute.PATCH("/reset-password", nil)
	// userRoute.PATCH("/change-pin",nil)
	userRoute.PATCH("/change-pin", middlewares.CheckAuthentication, middlewares.CheckAccountActivation, handlers.ChangeUserPin)
	userRoute.PATCH("/change-password", middlewares.CheckAuthentication, middlewares.CheckAccountActivation, handlers.ChangeUserPassword)
	// userRoute.PATCH("/change-email", middlewares.CheckAuthentication, middlewares.CheckAccountActivation, nil)
	userRoute.POST("/toggle-account-deactivation", middlewares.CheckAuthentication, handlers.ToogleAccountDeactivation)
}
