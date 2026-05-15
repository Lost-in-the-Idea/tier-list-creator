package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"tierlist/database/models"
	"time"
)

func HandleDatabaseActions(db *Database) error {
	actions := os.Getenv("DB_ACTION")
	env := os.Getenv("APP_ENV")
	if env != "dev" && (strings.Contains(actions, "clear") || strings.Contains(actions, "seed")) {
		return fmt.Errorf("Clear and Seed actions are only allowed in dev environment")
	}
	if actions == "" {
		return nil
	}

	actionList := strings.Split(actions, ",")
	for _, action := range actionList {
		switch strings.TrimSpace(action) {
		case "migrate":
			err := migrateDatabase(db)
			if err != nil {
				return fmt.Errorf("Failed to migrate database: %v", err)
			}
		case "seed":
			err := seedDatabase(db)
			if err != nil {
				return fmt.Errorf("Failed to seed database: %v", err)
			}
		case "clear":
			err := clearDatabase(db)
			if err != nil {
				return fmt.Errorf("Failed to clear database: %v", err)
			}
		default:
			fmt.Printf("Unknown database action: %s\n", action)
		}
	}
	os.Exit(0)
	return nil
}

func clearDatabase(db *Database) error {
	if db.DB == nil {
		return fmt.Errorf("Database connection unavailable")
	}
	err := db.DB.Exec("TRUNCATE TABLE users, tierlists, tierlist_items, submissions, submission_rankings, sessions RESTART IDENTITY CASCADE").Error
	if err != nil {
		return fmt.Errorf("Failed to drop tables: %v", err)
	}

	fmt.Println("Database Cleared Successfully")
	return nil
}

func migrateDatabase(db *Database) error {
	if db.DB == nil {
		return fmt.Errorf("Database connection unavailable")
	}
	err := db.DB.Migrator().DropTable(&models.User{}, &models.Tierlist{}, &models.TierlistItem{}, &models.Submissions{}, &models.SubmissionRankings{}, &models.Session{})
	if err != nil {
		return fmt.Errorf("Failed to drop tables: %v", err)
	}
	err = db.DB.AutoMigrate(&models.User{}, &models.Tierlist{}, &models.TierlistItem{}, &models.Submissions{}, &models.SubmissionRankings{}, &models.Session{})
	if err != nil {
		return fmt.Errorf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database Migrated Successfully")
	return nil
}

func seedDatabase(db *Database) error {
	if db.DB == nil {
		return fmt.Errorf("Database connection unavailable")
	}
	err := seedUsers(db)
	if err != nil {
		return fmt.Errorf("Failed to seed users: %v", err)
	}

	err = seedTierlists(db)
	if err != nil {
		return fmt.Errorf("Failed to seed tierlists: %v", err)
	}
	return nil
}

func seedUsers(db *Database) error {
	var usersJson []byte
	usersJson, err := os.ReadFile("database/seeds/users.json")
	if err != nil {
		return fmt.Errorf("Failed to read users seed file: %v", err)
	}

	var count int64
	db.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Println("Database already seeded, skipping seeding.")
		return nil
	}

	var users []models.User
	if err := json.Unmarshal(usersJson, &users); err != nil {
		return fmt.Errorf("Failed to unmarshal users seed data: %v", err)
	}
	for _, user := range users {
		user.LastLogin = time.Now()
		if err := db.DB.Create(&user).Error; err != nil {
			return fmt.Errorf("Failed to seed users: %v", err)
		}
	}
	fmt.Println("Users Seeded Successfully")
	return nil
}

func seedTierlists(db *Database) error {
	var tierlistsJson []byte
	tierlistsJson, err := os.ReadFile("database/seeds/tierlists.json")
	if err != nil {
		return fmt.Errorf("Failed to read tierlists seed file: %v", err)
	}

	var count int64
    db.DB.Model(&models.Tierlist{}).Count(&count)
    if count > 0 {
        fmt.Println("Tierlists already seeded, skipping")
        return nil
    }

    var tierlists []models.Tierlist
    if err := json.Unmarshal(tierlistsJson, &tierlists); err != nil {
        return fmt.Errorf("Failed to parse tierlists seed: %w", err)
    }

    for _, tl := range tierlists {
        if err := db.DB.Create(&tl).Error; err != nil {
            return fmt.Errorf("Failed to seed tierlist: %w", err)
        }
    }

	fmt.Println("Tierlists Seeded Successfully")
	return nil
}