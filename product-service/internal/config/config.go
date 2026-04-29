package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/joho/godotenv"
)

// Config holds all configuration for product-service.
// Values are resolved in this priority order:
//  1. GCP Secret Manager  (if GCP_PROJECT_ID is set)
//  2. Environment variables
//  3. .env file
//  4. Hard-coded defaults
type Config struct {
	HTTPPort   string
	DbDSN      string
	GCPProject string // optional – triggers Secret Manager lookup
}

// Load resolves config from env / GCP Secret Manager.
func Load(ctx context.Context) (*Config, error) {
	// Try loading .env (ignored if not present)
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort:   getEnv("HTTP_PORT", "8083"),
		GCPProject: os.Getenv("GCP_PROJECT_ID"),
	}

	// If GCP_PROJECT_ID is set, pull secrets from Secret Manager.
	if cfg.GCPProject != "" {
		dsn, err := getSecret(ctx, cfg.GCPProject, "PRODUCT_DB_DSN")
		if err != nil {
			log.Printf("[config] Secret Manager error (falling back to env): %v", err)
			cfg.DbDSN = mustEnv("DB_DSN")
		} else {
			cfg.DbDSN = dsn
		}
	} else {
		cfg.DbDSN = mustEnv("DB_DSN")
	}

	return cfg, nil
}

// getSecret fetches the latest version of a secret from GCP Secret Manager.
func getSecret(ctx context.Context, projectID, secretName string) (string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("secretmanager.NewClient: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)
	req := &secretmanagerpb.AccessSecretVersionRequest{Name: name}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("AccessSecretVersion(%s): %w", name, err)
	}
	return string(result.Payload.Data), nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("[config] required env var %q is not set", key)
	}
	return v
}

// GetInt reads an int env var, returning def if absent or invalid.
func GetInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
