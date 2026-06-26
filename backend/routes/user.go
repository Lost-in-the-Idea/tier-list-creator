package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tierlist/database/models"
	"tierlist/services"
)

func SetupUserRoutes(api *gin.RouterGroup, svc *services.UserService, authRequired gin.HandlerFunc) {
	users := api.Group("/users")
	users.Use(authRequired)

	users.GET("/", func(c *gin.Context) { getAllUsers(c, svc) })
	users.GET("/:id", func(c *gin.Context) { getUserById(c, svc) })
	users.POST("/", func(c *gin.Context) { createUser(c, svc) })
	users.DELETE("/:id", func(c *gin.Context) { deleteUser(c, svc) })
}

func getAllUsers(c *gin.Context, svc *services.UserService) {
	users, err := svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func getUserById(c *gin.Context, svc *services.UserService) {
	id := c.Param("id")
	user, err := svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context, svc *services.UserService) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := svc.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User Created",
	})
}

func deleteUser(c *gin.Context, svc *services.UserService) {
	id := c.Param("id")
	if err := svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User Deleted " + id,
	})
}
