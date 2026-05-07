package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"tierlist/database"
	"tierlist/utilities"
)

func main() {
	db := database.Database{}
	_ = godotenv.Load()
	var DBName = os.Getenv("DB_NAME")
	var DBUser = os.Getenv("DB_USER")
	var DBPassword = os.Getenv("DB_PASSWORD")
	var DBHost = os.Getenv("DB_HOST")
	var DBPort = os.Getenv("DB_PORT")

	err := db.InitialiseDatabase(DBName, DBUser, DBPassword, DBHost, DBPort)
	if err != nil {
		panic(err)
	}
	err = database.HandleDatabaseActions(&db)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()
	utilities.SetupRoutes(r, &db)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
  }