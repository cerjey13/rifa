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
	db, err := database.Connect(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to start the db: %v", err)
	}
	dbAdapter := database.NewPgxpoolAdapter(db)
	server, err := core.NewHttpServer(dbAdapter, core.HttpServerOptions{
		Logger: logger,
		Host:   "localhost",
		Port:   cfg.Port,
	})
	if err != nil {
		log.Fatal("failed to config the server")
	}
	server.Handle("/*", http.FileServer(http.FS(dist)))

	log.Println("Rifa backend listening on :" + cfg.Port)
	log.Fatal(server.ListenAndServe())
}
