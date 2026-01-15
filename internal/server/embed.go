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
	dist, err := fs.Sub(FrontendFS, "frontend/dist")
	var root fs.FS
	if err != nil {
		// Fallback to local disk
		root = os.DirFS("./frontend/dist")
	} else {
		root = dist
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
