package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DiscordID string `json:"discord_id" gorm:"uniqueIndex"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	LastLogin time.Time `json:"last_login"`
	Sessions []Session `json:"sessions" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // association to the Session model
}

type Session struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Token string `json:"token" gorm:"uniqueIndex"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	LastUsed time.Time `json:"last_used"`
	ExpiresAt time.Time `json:"expires_at"`
	}