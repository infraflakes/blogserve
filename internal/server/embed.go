package server

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// FrontendFS is a global embedded filesystem containing the built frontend.
// It is expected to be populated by the main package.
var FrontendFS embed.FS

// getFrontendHandler returns an http.Handler that serves the frontend.
// It tries to serve from the embedded FrontendFS first, falling back to the local
// "frontend/dist" directory if the embedded FS is missing or empty.
// It includes SPA routing logic: if a file is not found, it serves index.html.
func (s *Server) getFrontendHandler() http.Handler {
	var root fs.FS

	if s.FrontendFS != nil {
		root = s.FrontendFS
	} else {
		// Try embedded FS
		dist, err := fs.Sub(FrontendFS, "frontend/dist")
		if err == nil {
			// Check if it's actually populated
			if _, err := fs.ReadDir(dist, "."); err == nil {
				root = dist
			}
		}

		if root == nil {
			// Fallback to local disk. Try a few common locations.
			paths := []string{"./frontend/dist", "../../frontend/dist", "../frontend/dist"}
			for _, p := range paths {
				if _, err := os.Stat(p); err == nil {
					root = os.DirFS(p)
					break
				}
			}
		}
	}

	if root == nil {
		// Final fallback (might still fail, but better than nil)
		root = os.DirFS("./frontend/dist")
	}

	fsServer := http.FileServer(http.FS(root))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the file doesn't exist in the FS, serve index.html
		path := r.URL.Path
		if path != "/" {
			// Clean path to prevent traversal
			path = filepath.Clean(path)
			// Remove leading slash for fs.Stat
			fPath := path[1:]
			_, err := fs.Stat(root, fPath)
			if err != nil {
				// File not found, serve index.html for SPA routing
				r.URL.Path = "/"
			}
		}
		fsServer.ServeHTTP(w, r)
	})
}
