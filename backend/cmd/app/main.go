package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"rifa/backend/internal/core"
	"rifa/backend/pkg/config"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/logger"

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
		log.Fatal("failed to load environment variables")
	}

	logger := logger.New(cfg.Env)
	driver := database.NewPostgresDriver()
	dbAdapter, err := database.Connect(context.Background(), driver, cfg)
	if err != nil {
		log.Fatalf("failed to start the db: %v", err)
	}
	defer dbAdapter.Close()

	front := http.FS(dist)
	server, err := core.NewHttpServer(
		dbAdapter,
		front,
		core.HttpServerOptions{
			Logger:        logger,
			Host:          cfg.Host,
			Port:          cfg.Port,
			SecureCookies: cfg.UseSecureCookie,
		},
	)
	if err != nil {
		log.Fatal("failed to config the server")
	}

	log.Println("Rifa backend listening on :" + cfg.Port)
	log.Fatal(server.ListenAndServe())
}
