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

type Config struct {
	HTTPPort   string
	JWTSecret  string
	AuthURL    string
	AccountURL string
	ProductURL string
	OrderURL   string
	GCPProject string
}

func Load(ctx context.Context) (*Config, error) {
	_ = godotenv.Load()
	gcpProject := os.Getenv("GCP_PROJECT_ID")
	jwtRaw := MustEnv("JWT_SECRET")
	jwtSecret, err := resolveSecret(ctx, gcpProject, jwtRaw)
	if err != nil {
		log.Printf("[config] JWT_SECRET resolve failed, using raw value: %v", err)
		jwtSecret = jwtRaw
	}
	return &Config{
		HTTPPort:   GetEnv("HTTP_PORT", "8080"),
		JWTSecret:  jwtSecret,
		GCPProject: gcpProject,
		AuthURL:    GetEnv("AUTH_SERVICE_URL", "http://auth-service:8081"),
		AccountURL: GetEnv("ACCOUNT_SERVICE_URL", "http://account-service:8082"),
		ProductURL: GetEnv("PRODUCT_SERVICE_URL", "http://product-service:8083"),
		OrderURL:   GetEnv("ORDER_SERVICE_URL", "http://order-service:8084"),
	}, nil
}

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}

func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" { log.Fatalf("[config] required env var %q not set", key) }
	return v
}

func resolveSecret(ctx context.Context, projectID, value string) (string, error) {
	const prefix = "gcp://"
	if projectID == "" || len(value) <= len(prefix) || value[:len(prefix)] != prefix {
		return value, nil
	}
	secretName := value[len(prefix):]
	client, err := secretmanager.NewClient(ctx)
	if err != nil { return "", err }
	defer client.Close()
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: name})
	if err != nil { return "", err }
	return string(result.Payload.Data), nil
}
