package routes

import (
	"net/http"

	"tierlist/database"
	"tierlist/models"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	users := r.Group("/users")

	users.GET("/", getAllUsers) 
	users.GET("/:id", getUserById)
	users.POST("/", createUser)
	users.DELETE("/:id", deleteUser)

}

func getAllUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H {"Error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func getUserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User with UserID" + id,
	})
}

func createUser(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message": "User Created",
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User Deleted" + id,
	})
}