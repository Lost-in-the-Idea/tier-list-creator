package authentication

import (
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}
}