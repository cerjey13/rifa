package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"rifa/backend/internal/core"
	"rifa/backend/pkg/config"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/logger"
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

	logger := logger.New(cfg.Server.Env)
	driver := database.NewPostgresDriver()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbAdapter, err := database.Connect(ctx, driver, &cfg.Database)
	if err != nil {
		log.Fatalf("failed to start the db: %v", err)
	}
	defer dbAdapter.Close()

	shutdown, err := telemetry.InitCollector(ctx, cfg.Collector)
	if err != nil {
		log.Fatalf("collector init: %v", err)
	}
	defer func() {
		err = shutdown(context.Background())
		if err != nil {
			log.Fatalf("failed to shutdown collector: %v", err)
		}
	}()

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
		log.Fatal("failed to config the server")
	}

	log.Println("Rifa backend listening on :" + cfg.Server.Port)
	log.Fatal(server.ListenAndServe())
}
