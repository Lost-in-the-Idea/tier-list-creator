package routes

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"tierlist/database"
	"tierlist/database/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// These functions may be moved to seperate concerns and move them from the routes package TBC

func generateStateCookie() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func exchangeCodeForToken(c *gin.Context, code string) (*oauth2.Token, error) {
	token, err := conf.Exchange(c, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getUserInfoFromDiscord(c *gin.Context, token *oauth2.Token) (map[string]interface{}, error) {
	client := conf.Client(c, token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discord API returned status %d", resp.StatusCode)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func findOrCreateUser(db *database.Database, userInfo map[string]interface{}) (*models.User, error) {
	var user models.User
	discordID := userInfo["id"].(string)
	result := db.DB.Where("discord_id = ?", discordID).First(&user)
	if result.Error != nil {
		avatar, _ := userInfo["avatar"].(string)
		user = models.User{
			DiscordID: discordID,
			Username:  userInfo["username"].(string),
			Avatar:    avatar,
			LastLogin: time.Now(),
		}
		if err := db.DB.Create(&user).Error; err != nil {
			return nil, err
		}
	} else {
		avatar, _ := userInfo["avatar"].(string)
		user.Username = userInfo["username"].(string)
		user.Avatar = avatar
		user.LastLogin = time.Now()
		if err := db.DB.Save(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func createSession(db *database.Database, user *models.User) (*models.Session, error) {
	sessionToken := uuid.New().String()
	session := models.Session{
		Token:    sessionToken,
		UserID:   user.ID,
		LastUsed: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 168),
	}
	if err := db.DB.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func deleteSession(db *database.Database, session models.Session) error {
	return db.DB.Delete(&session).Error
}

func DeleteExpiredSessions(db *database.Database) {
	result := db.DB.Where("expires_at < ?", time.Now()).Delete(&models.Session{})
	if result.Error != nil {
		fmt.Printf("Error deleting expired sessions: %v\n", result.Error)
		return
	}
	fmt.Printf("Deleted %d expired sessions\n", result.RowsAffected)
}