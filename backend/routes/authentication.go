package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"tierlist/database/models"
	"tierlist/dto"
	"tierlist/middleware"
	"tierlist/services"
)

func SetupAuthenticationRoutes(api *gin.RouterGroup, svc *services.AuthService, cookieDomain, frontendURL string) {
	authentication := api.Group("/auth")
	// Rate-limit the OAuth entry points to blunt automated abuse of the login flow.
	authLimit := middleware.RateLimit(20, time.Minute)
	authentication.GET("/discord/redirect", authLimit, func(c *gin.Context) { handleDiscordRedirect(c, svc, cookieDomain) })
	authentication.GET("/discord/callback", authLimit, middleware.ValidateAuthState(cookieDomain), func(c *gin.Context) { handleDiscordCallback(c, svc, cookieDomain, frontendURL) })

	protected := authentication.Group("/")
	protected.Use(middleware.AuthRequired(svc, cookieDomain))
	protected.GET("/logout", func(c *gin.Context) { handleLogout(c, svc, cookieDomain) })
	protected.GET("/me", func(c *gin.Context) { getCurrentUser(c) })
}

func handleDiscordRedirect(c *gin.Context, svc *services.AuthService, cookieDomain string) {
	state, err := svc.GenerateStateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	c.SetCookie("login_state", state, 300, "/", cookieDomain, true, true)
	c.Redirect(http.StatusTemporaryRedirect, svc.BuildAuthURL(state))
}

func handleDiscordCallback(c *gin.Context, svc *services.AuthService, cookieDomain, frontendURL string) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
		return
	}
	token, err := svc.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		log.Printf("oauth: exchange code failed: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Login failed, please try again"})
		return
	}
	userInfo, err := svc.GetDiscordUserInfo(c.Request.Context(), token)
	if err != nil {
		log.Printf("oauth: fetch discord user failed: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Login failed, please try again"})
		return
	}
	user, err := svc.FindOrCreateUser(userInfo)
	if err != nil {
		log.Printf("oauth: find/create user failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed, please try again"})
		return
	}
	session, err := svc.CreateSession(user)
	if err != nil {
		log.Printf("oauth: create session failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed, please try again"})
		return
	}
	c.SetCookie("session_token", session.Token, 60*60*24*7, "/", cookieDomain, true, true)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

func handleLogout(c *gin.Context, svc *services.AuthService, cookieDomain string) {
	sessionVal, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}
	if err := svc.DeleteSession(sessionVal.(models.Session)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}
	c.SetCookie("session_token", "", -1, "/", cookieDomain, true, true)
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
