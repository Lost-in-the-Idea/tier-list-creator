package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Config holds all runtime configuration, loaded from environment variables.
type Config struct {
	AppEnv       string
	Port         string
	AppURL       string // public base URL of the backend, used to build the OAuth redirect URL
	FrontendURL  string // where the OAuth callback redirects the browser after login
	CookieDomain string

	DiscordClientID     string
	DiscordClientSecret string

	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBSSLMode  string
}

// Load reads configuration from the environment and fails fast when anything
// required is missing or obviously misconfigured for the target environment.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:              getEnv("APP_ENV", "dev"),
		Port:                getEnv("PORT", "8080"),
		AppURL:              getEnv("APP_URL", "http://localhost:8080"),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:4200"),
		CookieDomain:        os.Getenv("COOKIE_DOMAIN"),
		DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		DBName:              os.Getenv("DB_NAME"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
	}

	required := map[string]string{
		"DISCORD_CLIENT_ID":     cfg.DiscordClientID,
		"DISCORD_CLIENT_SECRET": cfg.DiscordClientSecret,
		"DB_NAME":               cfg.DBName,
		"DB_USER":               cfg.DBUser,
		"DB_PASSWORD":           cfg.DBPassword,
		"DB_HOST":               cfg.DBHost,
	}
	var missing []string
	for key, val := range required {
		if strings.TrimSpace(val) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	// Production sanity checks. COOKIE_DOMAIN must be set (empty is the common
	// "forgot to configure it" mistake). `localhost` is allowed but warned about,
	// so local prod-parity testing works while real deploys get a nudge.
	if !cfg.IsDev() {
		if cfg.CookieDomain == "" {
			return nil, fmt.Errorf("COOKIE_DOMAIN must be set when APP_ENV != dev (use your real domain, or 'localhost' for local prod-parity testing)")
		}
		if cfg.CookieDomain == "localhost" {
			log.Println("WARNING: COOKIE_DOMAIN=localhost outside dev is only valid for local testing; set your real domain for a public deploy")
		}
		if !strings.HasPrefix(cfg.AppURL, "https://") {
			log.Println("WARNING: APP_URL is not HTTPS; Secure cookies require the app to be served over HTTPS")
		}
	}

	return cfg, nil
}

// IsDev reports whether the app is running in the development environment.
func (c *Config) IsDev() bool {
	return c.AppEnv == "dev"
}

// OAuthRedirectURL is the full Discord OAuth callback URL for this environment.
func (c *Config) OAuthRedirectURL() string {
	return strings.TrimRight(c.AppURL, "/") + "/api/auth/discord/callback"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); strings.TrimSpace(v) != "" {
		return v
	}
	return fallback
}
