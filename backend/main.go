package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"tierlist/database"
	"tierlist/routes"
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
	api := r.Group("/api")
	routes.SetupTierlistRoutes(api, &db)
	routes.SetupUserRoutes(api, &db)
	routes.SetupAuthenticationRoutes(api, &db)
	r.Run()
}