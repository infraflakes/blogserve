package server

import (
	"blogserve/internal/blog"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"
)

// Server represents the blog's web server.
type Server struct {
	BlogDir    string
	Port       int
	FrontendFS fs.FS
	reload     chan bool
}

// NewServer creates a new Server instance with the specified blog directory and port.
func NewServer(blogDir string, port int) *Server {
	return &Server{
		BlogDir: blogDir,
		Port:    port,
		reload:  make(chan bool),
	}
}

// Start starts the HTTP server and sets up the API and frontend handlers.
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/posts", s.handleGetPosts)
	mux.HandleFunc("/api/reload", s.handleReload)

	// Serve assets from the blog directory
	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir(s.BlogDir))))

	// Frontend handler
	feHandler := s.getFrontendHandler()

	// Catch-all handler for SPA and static files
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		feHandler.ServeHTTP(w, r)
	})

	// Wrap mux with logging middleware
	loggedMux := s.loggingMiddleware(mux)

	slog.Info("Server listening", "url", fmt.Sprintf("http://localhost:%d", s.Port))
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), loggedMux)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

// handleReload handles the Server-Sent Events (SSE) connection for live reloads.
func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-s.reload:
			_, _ = fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// TriggerReload sends a reload signal to all connected clients via SSE.
func (s *Server) TriggerReload() {
	select {
	case s.reload <- true:
	default:
		// Drop if no one is listening
	}
}

// handleGetPosts returns a list of all blog posts as JSON.
func (s *Server) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	posts, err := blog.ScanDirectory(s.BlogDir)
	if err != nil {
		slog.Error("Failed to scan directory for posts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(posts)
}
