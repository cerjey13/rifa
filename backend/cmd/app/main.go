package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:dist
var frontendAssets embed.FS

func main() {
	dist, err := fs.Sub(frontendAssets, "dist")
	if err != nil {
		log.Fatalf("Failed to locate embedded dist: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/api/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Go!"))
	}))

	// Static frontend
	mux.Handle("/", http.FileServer(http.FS(dist)))

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", mux)
}
