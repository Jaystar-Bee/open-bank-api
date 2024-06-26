package main

import (
	"net/http"

	"github.com/Jaystar-Bee/open-bank-api/db"
	doc "github.com/Jaystar-Bee/open-bank-api/docs"
	"github.com/Jaystar-Bee/open-bank-api/inits"
	"github.com/Jaystar-Bee/open-bank-api/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title						OPEN BANK API
//	@version					1.0
//	@description				OPEN BANK API
//	@contact.name				John Ayilara (Jaystar)
//	@contact.email				jbayilara@gmail.com
//	@contact.url				https://github.com/Jaystar-Bee
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@BasePath					/

var app *gin.Engine

func init() {
	_ = godotenv.Load()
	db.InitDatabase()

	// Cloudinary
	inits.InitCloudinary()
}

func main() {
	server := gin.Default()
	// DOCS

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type", "X-Requested-With", "Origin", "Accept"}
	server.Use(cors.New(config))
	doc.SwaggerInfo.BasePath = "/"
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// ROUTES
	routes.UserRoutes(server)
	routes.WalletRoutes(server)
	routes.TransactionRoutes(server)
	routes.UploadRoutes(server)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
