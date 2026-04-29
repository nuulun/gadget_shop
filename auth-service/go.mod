module auth-service

go 1.22

require (
	cloud.google.com/go/secretmanager v1.13.3
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/prometheus/client_golang v1.19.0
	golang.org/x/crypto v0.22.0
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.10
)
