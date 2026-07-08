package config

import (
	"testing"
)

// setEnv sets the full set of required vars via t.Setenv, then applies overrides.
func setRequired(t *testing.T, overrides map[string]string) {
	base := map[string]string{
		"DISCORD_CLIENT_ID":     "id",
		"DISCORD_CLIENT_SECRET": "secret",
		"DB_NAME":               "db",
		"DB_USER":               "user",
		"DB_PASSWORD":           "pw",
		"DB_HOST":               "localhost",
		"APP_ENV":               "dev",
	}
	for k, v := range base {
		if _, ok := overrides[k]; !ok {
			t.Setenv(k, v)
		}
	}
	for k, v := range overrides {
		t.Setenv(k, v)
	}
}

func TestLoad_MissingRequiredFails(t *testing.T) {
	setRequired(t, map[string]string{"DB_PASSWORD": ""})
	if _, err := Load(); err == nil {
		t.Fatal("expected error when DB_PASSWORD is missing, got nil")
	}
}

func TestLoad_DevDefaults(t *testing.T) {
	setRequired(t, nil)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OAuthRedirectURL() != "http://localhost:8080/api/auth/discord/callback" {
		t.Errorf("unexpected redirect URL: %s", cfg.OAuthRedirectURL())
	}
	if cfg.DBSSLMode != "disable" {
		t.Errorf("expected default sslmode=disable, got %s", cfg.DBSSLMode)
	}
}

func TestLoad_ProdRequiresCookieDomain(t *testing.T) {
	setRequired(t, map[string]string{"APP_ENV": "production", "COOKIE_DOMAIN": ""})
	if _, err := Load(); err == nil {
		t.Fatal("expected error when COOKIE_DOMAIN is empty in production, got nil")
	}
}

func TestLoad_ProdAllowsLocalhostForTesting(t *testing.T) {
	setRequired(t, map[string]string{
		"APP_ENV":       "production",
		"COOKIE_DOMAIN": "localhost",
		"APP_URL":       "https://localhost",
	})
	if _, err := Load(); err != nil {
		t.Fatalf("localhost cookie domain should be allowed (with a warning) for local testing: %v", err)
	}
}

func TestLoad_ProdValid(t *testing.T) {
	setRequired(t, map[string]string{
		"APP_ENV":       "production",
		"COOKIE_DOMAIN": "tierlist.example.com",
		"APP_URL":       "https://tierlist.example.com",
	})
	if _, err := Load(); err != nil {
		t.Fatalf("unexpected error for valid prod config: %v", err)
	}
}
