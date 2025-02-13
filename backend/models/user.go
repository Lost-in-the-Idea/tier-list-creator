package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    DiscordID    string    `json:"discord_id" gorm:"uniqueIndex"`
    Username     string    `json:"username"`
    Avatar       string    `json:"avatar"`
    LastLogin    time.Time `json:"last_login"`
    Sessions     []Session `json:"sessions" gorm:"foreignKey:DiscordID;references:DiscordID"`
}
