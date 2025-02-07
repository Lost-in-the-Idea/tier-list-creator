package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, gin.H{"user": userInfo})
}