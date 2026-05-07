package utilities

import (
	"tierlist/database"
	"tierlist/routes"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine, db *database.Database) {
	api := r.Group("/api")
	routes.SetupTierlistRoutes(api, db)
	routes.SetupUserRoutes(api, db)
	routes.SetupAuthenticationRoutes(api, db)
	
}