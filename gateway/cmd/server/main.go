package main

import (
	"context"
	"gateway/internal/config"
	"gateway/internal/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx); err != nil {
		log.Fatalf("[gateway] %v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load(ctx)
	if err != nil {
		return err
	}
	log.Printf("[gateway] config loaded — auth=%s account=%s product=%s order=%s",
		cfg.AuthURL, cfg.AccountURL, cfg.ProductURL, cfg.OrderURL)

	gw := handler.New(cfg)
	mux := http.NewServeMux()
	gw.Register(mux)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      handler.Metrics(handler.Logging(handler.CORS(mux))),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("[gateway] listening on :%s", cfg.HTTPPort)
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
