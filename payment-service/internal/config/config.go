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

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("[config] required env var %q is not set", key)
	}
	return v
}

func Load() {
	_ = godotenv.Load()
}

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
