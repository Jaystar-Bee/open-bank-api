package main

import (
	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db.InitDatabase()
}

func main() {
	server := gin.Default()
	routes.UserRoutes(server)
	routes.WalletRoutes(server)
	server.Run()
}
