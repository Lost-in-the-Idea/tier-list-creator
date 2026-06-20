package routes

import (
	"fmt"
	"net/http"
	"os"
	"tierlist/database"
	"tierlist/database/models"
	"tierlist/dto"
	"tierlist/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var conf = &oauth2.Config{}

func init() {
		_ = godotenv.Load()

		conf = &oauth2.Config{
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/api/auth/discord/callback",
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
		}
	}

func SetupAuthenticationRoutes(api *gin.RouterGroup, db *database.Database) {
	authentication := api.Group("/auth")
	authentication.GET("/discord/redirect", func(c *gin.Context) { handleDiscordRedirect(c) })
	authentication.GET("/discord/callback", middleware.ValidateAuthState, func(c *gin.Context) { handleDiscordCallback(c, db) })

	protected := authentication.Group("/")
	protected.Use(middleware.AuthRequired(db))
	protected.GET("/logout", func(c *gin.Context) { handleLogout(c, db) })
	protected.GET("/me", func(c *gin.Context) { getCurrentUser(c) })
}

func handleDiscordRedirect(c *gin.Context) {
	state, err := generateStateCookie()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	fmt.Printf("Redirecting to Discord with URL: %s\n", conf.RedirectURL)
	c.SetCookie("login_state", state, 300, "/", "localhost", true, true)
	url := conf.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleDiscordCallback(c *gin.Context, db *database.Database) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
		return
	}

	token, err := exchangeCodeForToken(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := getUserInfoFromDiscord(c, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := findOrCreateUser(db, userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session, err := createSession(db, user)
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

func handleLogout(c *gin.Context, db *database.Database) {
	sessionToken, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}

	err := deleteSession(db, sessionToken.(models.Session))
	if err != nil {
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