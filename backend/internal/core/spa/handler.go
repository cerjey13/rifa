package spa

import (
	"net/http"
	"path/filepath"
	"strings"
)

func SpaHandler(staticFS http.FileSystem) http.HandlerFunc {
	fileServer := http.FileServer(staticFS)
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		f, err := staticFS.Open(path)
		if err == nil {
			f.Close()
			setCacheHeaders(w, path)
			setSecurityHeaders(w)
			fileServer.ServeHTTP(w, r)
			return
		}
		// fallback to index.html
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	}
}

func setCacheHeaders(w http.ResponseWriter, path string) {
	path = strings.TrimPrefix(path, "/")

	switch {
	case isLongCacheAsset(path):
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	case strings.HasSuffix(path, ".html"):
		w.Header().Set("Cache-Control", "no-cache")

	default:
		w.Header().Set("Cache-Control", "public, max-age=86400")
	}
}

func isLongCacheAsset(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return strings.HasPrefix(path, "assets/") ||
		strings.HasPrefix(path, "fonts/") ||
		mapsToKnownAssetType(ext)
}

func mapsToKnownAssetType(ext string) bool {
	switch ext {
	case ".js", ".css", ".jpg", ".jpeg", ".png", ".webp", ".svg", ".woff", ".woff2", ".ttf", ".otf":
		return true
	default:
		return false
	}
}

func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set(
		"Content-Security-Policy",
		"default-src 'self'; img-src 'self' https: data:; script-src 'self'; style-src 'self' 'unsafe-inline'",
	)
	w.Header().Set(
		"Strict-Transport-Security",
		"max-age=63072000; includeSubDomains; preload",
	)
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("X-Frame-Options", "DENY")
}
