package config

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/joho/godotenv"
)

// GetEnv returns the env var value or def if not set.
func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// MustEnv returns the env var value or fatals.
func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("[config] required env var %q is not set", key)
	}
	return v
}

// Load reads .env file (if present) then environment.
// If GCP_PROJECT_ID is set, secrets are fetched from Secret Manager for
// any key whose value starts with "gcp://".
// Example: DB_DSN=gcp://MY_SECRET_NAME  →  resolved from Secret Manager.
func Load() {
	_ = godotenv.Load()
}

// ResolveSecret resolves a value: if projectID is non-empty and the value
// looks like a GCP secret reference ("gcp://SECRET_NAME"), it fetches the
// latest version from Secret Manager. Otherwise returns the value as-is.
func ResolveSecret(ctx context.Context, projectID, value string) (string, error) {
	if projectID == "" {
		return value, nil
	}
	const prefix = "gcp://"
	if len(value) <= len(prefix) || value[:len(prefix)] != prefix {
		return value, nil
	}
	secretName := value[len(prefix):]
	return GetSecretVersion(ctx, projectID, secretName)
}

// GetSecretVersion fetches the latest version of a GCP Secret Manager secret.
func GetSecretVersion(ctx context.Context, projectID, secretName string) (string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("secretmanager.NewClient: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: name})
	if err != nil {
		return "", fmt.Errorf("AccessSecretVersion(%s): %w", name, err)
	}
	return string(result.Payload.Data), nil
}
