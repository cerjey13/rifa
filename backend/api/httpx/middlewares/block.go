package mymiddlewares

import (
	"net/http"
	"strings"
)

// BlockScanners denies requests that match known exploit or scan paths.
func BlockScanners(next http.Handler) http.Handler {
	blocked := []string{
		".git", "wp-admin", "wp-includes", "xmlrpc.php",
		"wordpress", "cms", "shop", "test", "config",
		".env", ".bak", "setup-config.php",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.ToLower(r.URL.Path)
		for _, b := range blocked {
			if strings.Contains(path, b) {
				http.NotFound(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
