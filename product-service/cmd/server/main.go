package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"product-service/internal/config"
	"product-service/internal/handler"
	"product-service/internal/repository"
	"product-service/internal/seed"
	"product-service/internal/service"
	"product-service/migrations"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalf("[main] fatal: %v", err)
	}
}

func run(ctx context.Context) error {
	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load(ctx)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	log.Printf("[main] config loaded (port=%s, gcp_project=%q)", cfg.HTTPPort, cfg.GCPProject)

	// ── DB connection (with retry) ────────────────────────────────────────────
	sqlDB, gormDB, err := connectDB(ctx, cfg.DbDSN)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	log.Println("[main] connected to database ✓")

	// ── Migrations ────────────────────────────────────────────────────────────
	if err := migrations.Run(ctx, sqlDB); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}
	log.Println("[main] migrations complete ✓")

	// ── Repository / Service / Handler ────────────────────────────────────────
	repo := repository.New(gormDB)
	svc := service.New(repo)
	h := handler.New(svc)

	// ── Seed (runs only if table is empty) ────────────────────────────────────
	if err := seed.Run(ctx, repo); err != nil {
		// Non-fatal: log and continue so the service still starts.
		log.Printf("[main] seed warning: %v", err)
	}

	// ── HTTP server ───────────────────────────────────────────────────────────
	mux := http.NewServeMux()
	h.Register(mux)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      handler.Metrics(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start in background; block until ctx is cancelled.
	errCh := make(chan error, 1)
	go func() {
		log.Printf("[main] product-service listening on :%s", cfg.HTTPPort)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("[main] shutting down gracefully …")
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutCancel()
		return srv.Shutdown(shutCtx)
	case err := <-errCh:
		return err
	}
}

// ─── DB helpers ───────────────────────────────────────────────────────────────

func connectDB(ctx context.Context, dsn string) (*sql.DB, *gorm.DB, error) {
	const maxAttempts = 30
	const retryDelay = 2 * time.Second

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		},
	)

	var (
		sqlDB  *sql.DB
		gormDB *gorm.DB
		err    error
	)

	for i := 1; i <= maxAttempts; i++ {
		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
		if err == nil {
			sqlDB, err = gormDB.DB()
			if err == nil {
				if pingErr := sqlDB.PingContext(ctx); pingErr == nil {
					return sqlDB, gormDB, nil
				}
			}
		}
		log.Printf("[main] waiting for database … (%d/%d)", i, maxAttempts)
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(retryDelay):
		}
	}
	return nil, nil, fmt.Errorf("could not connect after %d attempts: %w", maxAttempts, err)
}
