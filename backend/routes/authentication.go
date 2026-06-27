package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tierlist/database/models"
	"tierlist/dto"
	"tierlist/middleware"
	"tierlist/services"
)

func SetupAuthenticationRoutes(api *gin.RouterGroup, svc *services.AuthService) {
	authentication := api.Group("/auth")
	authentication.GET("/discord/redirect", func(c *gin.Context) { handleDiscordRedirect(c, svc) })
	authentication.GET("/discord/callback", middleware.ValidateAuthState, func(c *gin.Context) { handleDiscordCallback(c, svc) })

	protected := authentication.Group("/")
	protected.Use(middleware.AuthRequired(svc))
	protected.GET("/logout", func(c *gin.Context) { handleLogout(c, svc) })
	protected.GET("/me", getCurrentUser)
}

func handleDiscordRedirect(c *gin.Context, svc *services.AuthService) {
	state, err := svc.GenerateStateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	c.SetCookie("login_state", state, 300, "/", "localhost", true, true)
	c.Redirect(http.StatusTemporaryRedirect, svc.BuildAuthURL(state))
}

func handleDiscordCallback(c *gin.Context, svc *services.AuthService) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
		return
	}
	token, err := svc.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := svc.GetDiscordUserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, err := svc.FindOrCreateUser(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	session, err := svc.CreateSession(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("session_token", session.Token, 60*60*24*7, "/", "localhost", true, true)
	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID.String(),
		DiscordID: user.DiscordID,
		Username:  user.Username,
		Avatar:    user.Avatar,
	})
}

func handleLogout(c *gin.Context, svc *services.AuthService) {
	sessionVal, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}
	if err := svc.DeleteSession(sessionVal.(models.Session)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}
	c.SetCookie("session_token", "", -1, "/", "localhost", true, true)
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Logged out successfully"})
}

func getCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	currentUser := user.(models.User)
	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        currentUser.ID.String(),
		DiscordID: currentUser.DiscordID,
		Username:  currentUser.Username,
		Avatar:    currentUser.Avatar,
	})
}
