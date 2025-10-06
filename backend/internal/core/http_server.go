package core

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"rifa/backend/api"
	"rifa/backend/internal/core/spa"
	"rifa/backend/pkg/config"
	"rifa/backend/pkg/db"
	"rifa/backend/pkg/logx"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimdw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/riandyrn/otelchi"
)

// FS is an interface that abstracts filesystem operations.
// It can be implemented by both real and embedded filesystems.
type FS interface {
	// Open opens the named file for reading
	Open(name string) (fs.File, error)

	// ReadDir reads the contents of the directory and returns a slice of directory entries
	ReadDir(name string) ([]fs.DirEntry, error)

	// ReadFile reads the named file and returns its contents
	ReadFile(name string) ([]byte, error)
}

// HttpServerOptions contains configuration options for the HTTP server
type HttpServerOptions struct {
	// ApiDocs is a list of API documentation to be served
	ApiDocs []huma.API

	// Mode specifies the application mode (debug/release)
	// Mode types.Mode

	// I18nBundle is used for internationalization
	// I18nBundle *I18nBundle

	// Logger is used for server-related logging
	Logger logx.Logger

	// ServerOpts specifies the server config options
	ServerOpts config.ServerOpts

	// SpaFS is the embedded filesystem for SPA files
	SpaFS FS

	// SpaDir is the path to the SPA files in the embedded filesystem
	SpaDir string

	// Middlewares is a list of middlewares to be applied to the server
	Middlewares []func(http.Handler) http.Handler

	// ServiceOpts have the environment variables to initialize services
	ServiceOpts config.ServiceOpts
}
type HttpServer struct {
	*chi.Mux
	*http.Server
	APIDocs []huma.API
	logger  logx.Logger
}

func NewHttpServer(
	db db.DB,
	front http.FileSystem,
	opts HttpServerOptions,
) (*HttpServer, error) {
	if opts.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	router := chi.NewRouter()
	router.Use(chimdw.Logger)
	router.Use(chimdw.RequestID)
	router.Use(chimdw.RealIP)
	router.Use(chimdw.Recoverer)
	router.Use(chimdw.Timeout(15 * time.Second))
	router.Use(httprate.LimitAll(50, 1*time.Second))
	router.Use(chimdw.Compress(4))
	router.Use(otelchi.Middleware("rifa", otelchi.WithChiRoutes(router)))

	apiConfig := huma.DefaultConfig("rifa", "1.0.0")
	apiConfig.CreateHooks = nil
	humaApi := humachi.New(router, apiConfig)
	api.RegisterHttpRoutes(humaApi, db, opts.Logger, opts.ServiceOpts)

	router.Get("/", spa.SpaHandler(front))
	router.NotFound(spa.SpaHandler(front))

	server := &HttpServer{
		router,
		&http.Server{
			Addr:              opts.ServerOpts.Host + ":" + opts.ServerOpts.Port,
			Handler:           router,
			ReadTimeout:       opts.ServerOpts.TimeOuts.Read,
			WriteTimeout:      opts.ServerOpts.TimeOuts.Write,
			ReadHeaderTimeout: opts.ServerOpts.TimeOuts.ReadHeader,
			IdleTimeout:       opts.ServerOpts.TimeOuts.Idle,
		},
		opts.ApiDocs,
		opts.Logger,
	}

	return server, nil
}
