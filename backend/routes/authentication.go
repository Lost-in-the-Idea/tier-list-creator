package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tierlist/database"
	"tierlist/database/models"
	"tierlist/dto"
	"tierlist/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	authentication.GET("/discord/callback", func(c *gin.Context) { handleDiscordCallback(c, db) })

	protected := authentication.Group("/")
	protected.Use(middleware.AuthRequired(db))
	protected.GET("/logout", func(c *gin.Context) { handleLogout(c, db) })
	protected.GET("/me", func(c *gin.Context) { getCurrentUser(c) })
}

func handleDiscordRedirect(c *gin.Context) {
	state, err := GenerateStateCookie()
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
	
	// verify state from cookie and query parameter
	loginState, err := c.Cookie("login_state")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No login state provided"})
		return
	}

	if loginState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login state"})
		return
	}

	c.SetCookie("login_state", "", -1, "/", "localhost", true, true)

	// exchange code for access token
	tok, err := conf.Exchange(c, code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// use access token to get user info from Discord API
	client := conf.Client(c, tok)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// parse user info from response
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// check if user exists in database, if not create new user
	var user models.User
	discordID := userInfo["id"].(string)
	result := db.DB.Where("discord_id = ?", discordID).First(&user)
	if result.Error != nil {
		user = models.User{
			DiscordID: discordID,
			Username: userInfo["username"].(string),
			Avatar: userInfo["avatar"].(string),
			LastLogin: time.Now(),
		}
		if err := db.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		user.Username = userInfo["username"].(string)
		user.Avatar = userInfo["avatar"].(string)
		user.LastLogin = time.Now()
		if err := db.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// create session for user and set cookie
	sessionToken := uuid.New().String()
	session := models.Session{
		Token: sessionToken,
		UserID: user.ID,
		LastUsed: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 168),
	}

	if err := db.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// set secure flag to true in production
	c.SetCookie("session_token", sessionToken, 60*60*24*7, "/", "localhost", true, true)

	// redirect back to the frontend now that the session cookie is set
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4200"
	}
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

func handleLogout(c *gin.Context, db *database.Database) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No session token provided"})
		return
	}

	var session models.Session
	result := db.DB.Where("token = ?", sessionToken).First(&session)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session not found"})
		return
	}

	if err := db.DB.Delete(&session).Error; err != nil {
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