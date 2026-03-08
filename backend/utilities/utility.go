package utilities

import (
	"tierlist/routes"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	routes.SetupTierlistRoutes(api)
	routes.SetupUserRoutes(api)
	routes.SetupAuthenticationRoutes(api)
	
}