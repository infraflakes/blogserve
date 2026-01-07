package server

import (
	"blogserve/backend/internal/blog"
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
	// Each post directory has an "assets" sub-directory
	// We want to serve /assets/post-slug/image.png -> blogDir/post-slug/assets/image.png
	// For simplicity, we'll just serve the whole blogDir under /data/ and let frontend handle paths,
	// OR we can make a custom handler. Let's start with a generic /data/ for now.
	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir(s.BlogDir))))

	// Static files and frontend will be added later
	// For now, let's assume we serve them from a 'dist' folder if it exists
	// mux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

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
