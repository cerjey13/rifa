package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rifa/backend/internal/background"
	"rifa/backend/internal/core"
	"rifa/backend/pkg/config"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/logx"
	"rifa/backend/pkg/telemetry"

	_ "github.com/joho/godotenv/autoload"
)

//go:embed all:dist
var frontendAssets embed.FS

func main() {
	dist, err := fs.Sub(frontendAssets, "dist")
	if err != nil {
		log.Fatalf("Failed to locate embedded dist: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	driver := database.NewPostgresDriver()
	ctx := context.Background()
	dbCtx, cancelDb := context.WithTimeout(ctx, 5*time.Second)
	defer cancelDb()

	dbAdapter, err := database.Connect(dbCtx, driver, &cfg.Database)
	if err != nil {
		log.Fatalf("failed to start the db: %v", err)
	}
	defer dbAdapter.Close()

	shutdown, err := telemetry.InitCollector(ctx, cfg.Collector)
	if err != nil {
		log.Fatalf("collector init: %v", err)
	}
	defer func() {
		err = shutdown(ctx)
		if err != nil {
			log.Fatalf("failed to shutdown collector: %v", err)
		}
	}()

	logger := logx.NewLogger(cfg.Server.Env)
	signalCtx, cancelSignal := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancelSignal()

	background.StartIdempotencyCleanup(signalCtx, dbAdapter, logger)

	front := http.FS(dist)
	server, err := core.NewHttpServer(
		dbAdapter,
		front,
		core.HttpServerOptions{
			Logger:      logger,
			ServerOpts:  cfg.Server,
			ServiceOpts: cfg.Service,
		},
	)
	if err != nil {
		logger.Error("Failed to configure server", "err", err)
		os.Exit(1)
	}

	go func() {
		logger.Info("Rifa backend listening", "port", cfg.Server.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "err", err)
			cancelSignal()
		}
	}()

	<-signalCtx.Done()
	logger.Info("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "err", err)
	}

	logger.Info("Server exited cleanly")
}
