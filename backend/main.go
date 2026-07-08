package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	discord "github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"

	"tierlist/config"
	"tierlist/database"
	"tierlist/routes"
	"tierlist/services"

	"tierlist/middleware"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	db, err := database.NewDatabase(
		cfg.DBName, cfg.DBUser, cfg.DBPassword,
		cfg.DBHost, cfg.DBPort, cfg.DBSSLMode,
	)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	if err := database.HandleDatabaseActions(db.DB); err != nil {
		log.Fatalf("database action error: %v", err)
	}
	defer db.Close()

	oauthConf := &oauth2.Config{
		ClientID:     cfg.DiscordClientID,
		ClientSecret: cfg.DiscordClientSecret,
		RedirectURL:  cfg.OAuthRedirectURL(),
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}
	authSvc := services.NewAuthService(db.DB, oauthConf)
	tierlistSvc := services.NewTierlistService(db.DB)

	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			authSvc.DeleteExpiredSessions()
		}
	}()

	authRequired := middleware.AuthRequired(authSvc, cfg.CookieDomain)
	optionalAuth := middleware.OptionalAuth(authSvc, cfg.CookieDomain)

	if !cfg.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	// Default every cookie this app sets to SameSite=Lax. The frontend is served
	// same-origin behind the reverse proxy, so Lax keeps the top-level Discord
	// login redirect working while blocking cross-site cookie sending.
	r.Use(func(c *gin.Context) {
		c.SetSameSite(http.SameSiteLaxMode)
		c.Next()
	})

	// Lightweight liveness endpoint for container healthchecks and the proxy.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	routes.SetupTierlistRoutes(api, tierlistSvc, authRequired, optionalAuth)
	routes.SetupAuthenticationRoutes(api, authSvc, cfg.CookieDomain, cfg.FrontendURL)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
