package main

import (
	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	db.InitDatabase()
}

func main() {
	server := gin.Default()
	routes.UserRoutes(server)
	server.Run()
}
