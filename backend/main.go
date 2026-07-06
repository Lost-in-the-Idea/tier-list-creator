package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	discord "github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"

	"tierlist/database"
	"tierlist/middleware"
	"tierlist/routes"
	"tierlist/services"
)

func main() {
	_ = godotenv.Load()
	var cookieDomain = os.Getenv("COOKIE_DOMAIN")

	db, err := database.NewDatabase(
		os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
	)
	if err != nil {
		panic(err)
	}
	err = database.HandleDatabaseActions(db.DB)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	oauthConf := &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/api/auth/discord/callback",
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}
	authSvc := services.NewAuthService(db.DB, oauthConf)
	tierlistSvc := services.NewTierlistService(db.DB)
	userSvc := services.NewUserService(db.DB)

	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			authSvc.DeleteExpiredSessions()
		}
	}()

	authRequired := middleware.AuthRequired(authSvc, cookieDomain)
	optionalAuth := middleware.OptionalAuth(authSvc, cookieDomain)

	r := gin.Default()
	api := r.Group("/api")
	routes.SetupTierlistRoutes(api, tierlistSvc, authRequired, optionalAuth)
	routes.SetupUserRoutes(api, userSvc, authRequired)
	routes.SetupAuthenticationRoutes(api, authSvc, cookieDomain)
	r.Run()
}
