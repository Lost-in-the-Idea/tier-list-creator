package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tierlist/services"
)

func SetupUserRoutes(api *gin.RouterGroup, svc *services.UserService, authRequired gin.HandlerFunc) {
	users := api.Group("/users")
	users.Use(authRequired)

	users.GET("/", func(c *gin.Context) { getAllUsers(c, svc) })
	users.GET("/:id", getUserById)
	users.POST("/", createUser)
	users.DELETE("/:id", deleteUser)
}

func getAllUsers(c *gin.Context, svc *services.UserService) {
	users, err := svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func getUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User with UserID " + id,
	})
}

func createUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User Created",
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User Deleted " + id,
	})
}
