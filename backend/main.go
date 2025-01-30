package main

import (
	"github.com/gin-gonic/gin"

	"tierlist/database"
	"tierlist/utilities"
)

func main() {
	database.ConnectDatabase()
	database.Seed()
	r := gin.Default()
	utilities.SetupRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
  }