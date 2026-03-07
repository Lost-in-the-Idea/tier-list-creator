package database

import (
	"database/sql"
	"fmt"
	"tierlist/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func ConnectDatabase() {
	sqlDB, err := sql.Open("sqlite", "mainframe.db")
	if err != nil {
		panic("Failed to open database with modernc driver: " + err.Error())
	}

	db, err := gorm.Open(sqlite.New(sqlite.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = db.AutoMigrate(&models.User{}, &models.Tierlist{}, &models.Tier{}, &models.Item{}, &models.Session{})
	if err != nil {
		panic("Failed to migrate database")
	}

	DB = db

	fmt.Print("Connection successful")
}