package routes

import (
	"net/http"

	"tierlist/database"
	"tierlist/database/models"
	"tierlist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup, db *database.Database) {
	users := api.Group("/users")
	users.Use(middleware.AuthRequired(db))

	users.GET("/", func(c *gin.Context) { getAllUsers(c, db) })
	users.GET("/:id", func(c *gin.Context) { getUserById(c, db) })
	users.POST("/", func(c *gin.Context) { createUser(c, db) })
	users.DELETE("/:id", func(c *gin.Context) { deleteUser(c, db) })

}

func getAllUsers(c *gin.Context, db *database.Database) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H {"Error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func getUserById(c *gin.Context, db *database.Database) { // not implemented yet
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User with UserID " + id,
	})
}

func createUser(c *gin.Context, db *database.Database){ // not implemented yet
	c.JSON(http.StatusOK, gin.H{
		"message": "User Created",
	})
}

func deleteUser(c *gin.Context, db *database.Database) { // not implemented yet
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User Deleted " + id,
	})
}