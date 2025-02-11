package middleware

import (
	"net/http"
	"tierlist/database"
	"tierlist/models"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func (c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var session models.Session
		if err := database.DB.Where("token = ?", sessionToken).First(&session).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if time.Now().After(session.ExpiresAt) {
			database.DB.Delete(&session)
			c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session Expired"})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.Where("discord_id = ?", session.DiscordID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not Found"})
			c.Abort()
			return
		}

		now := time.Now()
		if err := database.DB.Model(&session).Update("expires_at", now.Add(time.Hour * 168)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			c.Abort()
			return
		}
		if err := database.DB.Model(&session).Update("last_used", now).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			c.Abort()
			return
		}


		c.Set("user", user)
		c.Set("session", session)
		c.Next()
	}
}