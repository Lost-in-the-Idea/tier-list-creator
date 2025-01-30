package utilities

import (
	"tierlist/routes"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine) {
	routes.SetupTierlistRoutes(r)
	routes.SetupUserRoutes(r)
	
}