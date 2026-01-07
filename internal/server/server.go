package server

import (
	"blogserve/internal/blog"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	BlogDir string
	Port    int
	reload  chan bool
}

func NewServer(blogDir string, port int) *Server {
	return &Server{
		BlogDir: blogDir,
		Port:    port,
		reload:  make(chan bool),
	}
}

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
		// If it's a file that exists, serve it
		// Otherwise, serve index.html
		feHandler.ServeHTTP(w, r)
	})

	fmt.Printf("Server listening on http://localhost:%d\n", s.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
}

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
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) TriggerReload() {
	select {
	case s.reload <- true:
	default:
		// Drop if no one is listening
	}
}

func (s *Server) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	posts, err := blog.ScanDirectory(s.BlogDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
