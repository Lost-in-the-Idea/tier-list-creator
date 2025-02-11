package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
    gorm.Model
    Token      string    `json:"token" gorm:"uniqueIndex"`
    DiscordID  string    `json:"discord_id"`
    ExpiresAt  time.Time `json:"expires_at"`
}