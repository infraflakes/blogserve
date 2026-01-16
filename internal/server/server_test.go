package server

import (
	"blogserve/internal/blog"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestHandleGetPosts(t *testing.T) {
	tests := []struct {
		name          string
		posts         map[string]map[string]string // dir -> filename -> content
		expectedCount int
		expectedTitle string
	}{
		{
			name: "single valid post",
			posts: map[string]map[string]string{
				"post1": {
					"content.md": "# Hello",
					"meta.json":  `{"title": "Test Post"}`,
				},
			},
			expectedCount: 1,
			expectedTitle: "Test Post",
		},
		{
			name:          "empty directory",
			posts:         map[string]map[string]string{},
			expectedCount: 0,
		},
		{
			name: "multiple posts",
			posts: map[string]map[string]string{
				"post1": {
					"a.md":      "# Post 1",
					"meta.json": `{"title": "Title 1"}`,
				},
				"post2": {
					"b.md":      "# Post 2",
					"meta.json": `{"title": "Title 2"}`,
				},
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "servertest")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			for dir, files := range tt.posts {
				postPath := filepath.Join(tempDir, dir)
				os.MkdirAll(postPath, 0755)
				for name, content := range files {
					os.WriteFile(filepath.Join(postPath, name), []byte(content), 0644)
				}
			}

			s := NewServer(tempDir, 8080)
			req := httptest.NewRequest("GET", "/api/posts", nil)
			rr := httptest.NewRecorder()
			s.handleGetPosts(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected 200, got %d", rr.Code)
			}

			var posts []blog.Post
			json.NewDecoder(rr.Body).Decode(&posts)

			if len(posts) != tt.expectedCount {
				t.Errorf("expected %d posts, got %d", tt.expectedCount, len(posts))
			}

			if tt.expectedTitle != "" && len(posts) > 0 {
				if posts[0].Metadata.Title != tt.expectedTitle {
					t.Errorf("expected title %s, got %s", tt.expectedTitle, posts[0].Metadata.Title)
				}
			}
		})
	}
}

func TestHandleDataServing(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "datatest")
	defer os.RemoveAll(tempDir)

	testFile := "image.png"
	testContent := "dummy image content"
	postDir := filepath.Join(tempDir, "post1")
	os.MkdirAll(postDir, 0755)
	os.WriteFile(filepath.Join(postDir, testFile), []byte(testContent), 0644)

	s := NewServer(tempDir, 8080)
	mux := http.NewServeMux()
	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir(s.BlogDir))))

	req := httptest.NewRequest("GET", "/data/post1/image.png", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != testContent {
		t.Errorf("expected %s, got %s", testContent, rr.Body.String())
	}
}

func TestHandleReload(t *testing.T) {
	s := NewServer(".", 8080)

	req := httptest.NewRequest("GET", "/api/reload", nil)
	ctx, cancel := context.WithTimeout(req.Context(), 100*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Use a channel to signal that the handler has started
	started := make(chan bool)
	go func() {
		started <- true
		s.handleReload(rr, req)
	}()

	<-started
	// Give it a tiny bit of time to start the loop
	time.Sleep(10 * time.Millisecond)

	s.TriggerReload()

	// Wait for context to time out or reload to be processed
	<-ctx.Done()

	body := rr.Body.String()
	if !strings.Contains(body, "data: reload") {
		t.Errorf("expected reload data in response, got %s", body)
	}
}

func TestSPAHandler(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "frontendtest")
	defer os.RemoveAll(tempDir)

	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte("<html>index</html>"), 0644)
	os.WriteFile(filepath.Join(tempDir, "style.css"), []byte("body {}"), 0644)

	s := NewServer(".", 8080)
	// Inject the mock FS directly
	s.FrontendFS = os.DirFS(tempDir)
	handler := s.getFrontendHandler()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"existing file", "/style.css", http.StatusOK, "body {}"},
		{"SPA fallback root", "/blog/my-post", http.StatusOK, "<html>index</html>"},
		{"SPA fallback nested", "/a/b/c", http.StatusOK, "<html>index</html>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	s := NewServer(".", 8080)
	handler := s.loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTeapot {
		t.Errorf("middleware failed to pass through, got %d", rr.Code)
	}
}
