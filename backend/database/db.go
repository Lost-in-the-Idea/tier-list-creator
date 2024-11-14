package database

import (
	"fmt"
	"tierlist/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("mainframe.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	err = db.AutoMigrate(&models.User{}, &models.Tierlist{}, &models.Tiers{}, &models.Items{})
	if err != nil {
		panic("Failed to migrate database")
	}

	DB = db

	fmt.Print("Connection successful")
}