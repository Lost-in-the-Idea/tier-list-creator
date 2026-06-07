package middleware

import (
	"net/http"
	"tierlist/database"
	"tierlist/database/models"
	"time"

	"github.com/gin-gonic/gin"
)

func OptionalAuth(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.Next()
			return
		}

		var session models.Session
		if err := db.DB.Where("token = ?", sessionToken).First(&session).Error; err != nil {
			c.Next()
			return
		}

		if time.Now().After(session.ExpiresAt) {
			db.DB.Delete(&session)
			c.SetCookie("session_token", "", -1, "/", "localhost", true, true)
			c.Next()
			return
		}

		var user models.User
		if err := db.DB.Where("id = ?", session.UserID).First(&user).Error; err != nil {
			c.Next()
			return
		}

		c.Set("user", user)
		c.Set("session", session)
		c.Next()
	}
}

func AuthRequired(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var session models.Session
		if err := db.DB.Where("token = ?", sessionToken).First(&session).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if time.Now().After(session.ExpiresAt) {
			db.DB.Delete(&session)
			c.SetCookie("session_token", "", -1, "/", "localhost", true, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session Expired"})
			c.Abort()
			return
		}

		var user models.User
		if err := db.DB.Where("id = ?", session.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not Found"})
			c.Abort()
			return
		}

		now := time.Now()
		expiryDuration := time.Hour * 168
		if err := db.DB.Model(&session).Update("expires_at", now.Add(expiryDuration)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			c.Abort()
			return
		}
		if err := db.DB.Model(&session).Update("last_used", now).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			c.Abort()
			return
		}

		c.SetCookie("session_token", sessionToken, int(expiryDuration.Seconds()), "/", "localhost", true, true)
		c.Set("user", user)
		c.Set("session", session)
		c.Next()
	}
}