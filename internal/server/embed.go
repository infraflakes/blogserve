package server

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// Global embed FS
var FrontendFS embed.FS

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

