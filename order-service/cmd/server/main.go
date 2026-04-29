package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"order-service/internal/config"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/service"
	"order-service/migrations"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx); err != nil {
		log.Fatalf("[order] %v", err)
	}
}

func run(ctx context.Context) error {
	config.Load()
	gcpProject := os.Getenv("GCP_PROJECT_ID")

	dbDSN, err := config.ResolveSecret(ctx, gcpProject, config.MustEnv("DB_DSN"))
	if err != nil {
		return fmt.Errorf("resolve DB_DSN: %w", err)
	}
	productURL := config.GetEnv("PRODUCT_SERVICE_URL", "http://product-service:8083")
	port := config.GetEnv("HTTP_PORT", "8084")

	sqlDB, gormDB, err := connectDB(ctx, dbDSN)
	if err != nil {
		return err
	}
	if err := migrations.Run(ctx, sqlDB); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	repo := repository.New(gormDB)
	svc := service.New(repo, productURL)
	h := handler.New(svc)

	mux := http.NewServeMux()
	h.Register(mux)

	srv := &http.Server{Addr: ":" + port, Handler: handler.Metrics(mux), ReadTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second}
	errCh := make(chan error, 1)
	go func() {
		log.Printf("[order] listening on :%s", port)
		errCh <- srv.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutCtx)
	case err := <-errCh:
		return err
	}
}

func connectDB(ctx context.Context, dsn string) (*sql.DB, *gorm.DB, error) {
	for i := 1; i <= 30; i++ {
		gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := gormDB.DB()
			if sqlDB.PingContext(ctx) == nil {
				return sqlDB, gormDB, nil
			}
		}
		log.Printf("[order] waiting for db (%d/30)…", i)
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	return nil, nil, fmt.Errorf("db unavailable")
}
