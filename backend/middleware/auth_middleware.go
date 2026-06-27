package middleware

import (
	"net/http"
	"tierlist/services"
	"time"

	"github.com/gin-gonic/gin"
)

func OptionalAuth(svc *services.AuthService, cookieDomain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			c.Next()
			return
		}
		session, user, err := svc.ResolveSession(token)
		if err != nil {
			c.Next()
			return
		}
		c.Set("user", *user)
		c.Set("session", *session)
		c.Next()
	}
}

func AuthRequired(svc *services.AuthService, cookieDomain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		session, user, err := svc.ResolveSession(token)
		if err != nil {
			if err.Error() == "session expired" {
				c.SetCookie("session_token", "", -1, "/", cookieDomain, true, true)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session Expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
			c.Abort()
			return
		}
		if err := svc.RollSession(session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			c.Abort()
			return
		}
		c.SetCookie("session_token", session.Token, int((time.Hour * 168).Seconds()), "/", cookieDomain, true, true)
		c.Set("user", *user)
		c.Set("session", *session)
		c.Next()
	}
}

func ValidateAuthState(cookieDomain string) gin.HandlerFunc {
	return func(c *gin.Context) {
	loginState, err := c.Cookie("login_state")
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No login state provided"})
			c.Abort()
			return
		}

	if loginState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login state"})
		c.Abort()
		return
	}
	
	// clear cookie after validating to prevent reuse
	c.SetCookie("login_state", "", -1, "/", cookieDomain, true, true)
	c.Next()
}
}