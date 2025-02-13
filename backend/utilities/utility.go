package utilities

import (
	"tierlist/routes"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://127.0.0.1:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:          12 * time.Hour,
    }))

	routes.SetupTierlistRoutes(r)
	routes.SetupUserRoutes(r)
	routes.SetupAuthenticationRoutes(r)
	
}