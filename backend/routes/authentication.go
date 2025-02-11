package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tierlist/database"
	"tierlist/middleware"
	"tierlist/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var (
	state = "random"
	conf = &oauth2.Config{}
)

func init() {
		err := godotenv.Load()
		if err != nil {
			panic("Failed to load .env file")
		}
		
		conf = &oauth2.Config{
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/auth/discord/callback",
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
		}
	}


func SetupAuthenticationRoutes(r *gin.Engine) {
	authentication := r.Group("/auth")
	authentication.GET("/discord/redirect", handleDiscordRedirect)
	authentication.GET("/discord/callback", handleDiscordCallback)

	protected := authentication.Group("/")
	protected.Use(middleware.AuthRequired())
	protected.GET("/logout", handleLogout)
	protected.GET("/me", getCurrentUser)
}

func handleDiscordRedirect(c *gin.Context) {
	fmt.Printf("Redirecting to Discord with URL: %s\n", conf.RedirectURL)
	url := conf.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleDiscordCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
		return
	}

	tok, err := conf.Exchange(c, code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := conf.Client(c, tok)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	var user models.User
	discordID := userInfo["id"].(string)
	result := database.DB.Where("discord_id = ?", discordID).First(&user)
	if result.Error != nil {
		user = models.User{
			DiscordID: discordID,
			Username: userInfo["username"].(string),
			Avatar: userInfo["avatar"].(string),
			LastLogin: time.Now(),
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		user.Username = userInfo["username"].(string)
		user.Avatar = userInfo["avatar"].(string)
		user.LastLogin = time.Now()
		if err := database.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	sessionToken := uuid.New().String()
	session := models.Session{
		Token: sessionToken,
		DiscordID: user.DiscordID,
		LastUsed: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 168),
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// set secure flag to true in production
	c.SetCookie("session_token", sessionToken, 60*60*24*7, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"user": user, "session": session})

}

func handleLogout(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No session token provided"})
		return
	}

	var session models.Session
	result := database.DB.Where("token = ?", sessionToken).First(&session)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session not found"})
		return
	}

	if err := database.DB.Delete(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func getCurrentUser(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"user": user})
}