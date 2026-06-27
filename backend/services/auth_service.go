package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"tierlist/database/models"
)

type AuthService struct {
	db   *gorm.DB
	conf *oauth2.Config
}

func NewAuthService(db *gorm.DB, conf *oauth2.Config) *AuthService {
	return &AuthService{db: db, conf: conf}
}

func (s *AuthService) GenerateStateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *AuthService) BuildAuthURL(state string) string {
	return s.conf.AuthCodeURL(state)
}

func (s *AuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.conf.Exchange(ctx, code)
}

func (s *AuthService) GetDiscordUserInfo(ctx context.Context, token *oauth2.Token) (map[string]interface{}, error) {
	client := s.conf.Client(ctx, token)
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

func (s *AuthService) FindOrCreateUser(userInfo map[string]interface{}) (*models.User, error) {
	var user models.User
	discordID := userInfo["id"].(string)
	if err := s.db.Where("discord_id = ?", discordID).First(&user).Error; err != nil {
		avatar, _ := userInfo["avatar"].(string)
		user = models.User{
			DiscordID: discordID,
			Username:  userInfo["username"].(string),
			Avatar:    avatar,
			LastLogin: time.Now(),
		}
		if err := s.db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else {
		avatar, _ := userInfo["avatar"].(string)
		user.Username = userInfo["username"].(string)
		user.Avatar = avatar
		user.LastLogin = time.Now()
		if err := s.db.Save(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *AuthService) CreateSession(user *models.User) (*models.Session, error) {
	session := models.Session{
		Token:     uuid.New().String(),
		UserID:    user.ID,
		LastUsed:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 168),
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *AuthService) DeleteSession(session models.Session) error {
	return s.db.Delete(&session).Error
}

func (s *AuthService) RollSession(session *models.Session) error {
	now := time.Now()
	const expiryDuration = time.Hour * 168
	if err := s.db.Model(session).Updates(map[string]any{
		"expires_at": now.Add(expiryDuration),
		"last_used":  now,
	}).Error; err != nil {
		return err
	}
	session.ExpiresAt = now.Add(expiryDuration)
	session.LastUsed = now
	return nil
}

func (s *AuthService) DeleteExpiredSessions() {
	result := s.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{})
	if result.Error != nil {
		fmt.Printf("Error deleting expired sessions: %v\n", result.Error)
		return
	}
	fmt.Printf("Deleted %d expired sessions\n", result.RowsAffected)
}

func (s *AuthService) ResolveSession(c *gin.Context, cookieDomain string) (*models.Session, *models.User, error) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		return nil, nil, err
	}

	var session models.Session
	if err := s.db.Where("token = ?", sessionToken).First(&session).Error; err != nil {
		return nil, nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		s.db.Delete(&session)
		c.SetCookie("session_token", "", -1, "/", cookieDomain, true, true)
		return nil, nil, errors.New("session expired")
	}

	var user models.User
	if err := s.db.Where("id = ?", session.UserID).First(&user).Error; err != nil {
		return nil, nil, err
	}

	return &session, &user, nil
}
