package core

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"rifa/backend/api/httpx"
	"rifa/backend/pkg/db"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimdw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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
	Logger *slog.Logger

	// Host specifies the server binding address
	Host string

	// Port specifies the server listening port
	Port string

	// SpaFS is the embedded filesystem for SPA files
	SpaFS FS

	// SpaDir is the path to the SPA files in the embedded filesystem
	SpaDir string

	// Middlewares is a list of middlewares to be applied to the server
	Middlewares []func(http.Handler) http.Handler

	//
	SecureCookies bool
}
type HttpServer struct {
	*chi.Mux
	*http.Server
	APIDocs []huma.API
	logger  *slog.Logger
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
	router.Use(httprate.LimitAll(10, 1*time.Second))
	router.NotFound(spaHandler(front))
	router.Use(chimdw.Compress(4))
	router.Use(chimdw.Recoverer)

	apiConfig := huma.DefaultConfig("rifa", "1.0.0")
	apiConfig.CreateHooks = nil
	api := humachi.New(router, apiConfig)
	httpx.RegisterAuthRoutes(api, opts.SecureCookies, db)
	httpx.RegisterPurchaseRoutes(api, db)
	httpx.RegisterTicketsRoutes(api, db)

	server := &HttpServer{
		router,
		&http.Server{
			Addr:    opts.Host + ":" + opts.Port,
			Handler: router,
		},
		opts.ApiDocs,
		opts.Logger,
	}

	return server, nil
}

func spaHandler(staticFS http.FileSystem) http.HandlerFunc {
	fileServer := http.FileServer(staticFS)
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		f, err := staticFS.Open(path)
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		// fallback to index.html
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	}
}
