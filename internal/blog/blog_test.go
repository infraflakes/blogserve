package blog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadPost(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "blogtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("ValidPost", func(t *testing.T) {
		postDir := filepath.Join(tempDir, "valid-post")
		os.Mkdir(postDir, 0755)
		os.WriteFile(filepath.Join(postDir, "hello.md"), []byte("# Hello"), 0644)
		os.WriteFile(filepath.Join(postDir, "meta.json"), []byte(`{"title": "Meta"}`), 0644)

		post, err := readPost(postDir, "valid-post")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if post.Content != "# Hello" {
			t.Errorf("expected content # Hello, got %s", post.Content)
		}
		if post.Metadata.Title != "Meta" {
			t.Errorf("expected title Meta, got %s", post.Metadata.Title)
		}
	})

	t.Run("NoMarkdown", func(t *testing.T) {
		postDir := filepath.Join(tempDir, "no-md")
		os.Mkdir(postDir, 0755)

		_, err := readPost(postDir, "no-md")
		if err == nil {
			t.Error("expected error for missing markdown file")
		}
	})

	t.Run("MultipleMarkdown", func(t *testing.T) {
		postDir := filepath.Join(tempDir, "multi-md")
		os.Mkdir(postDir, 0755)
		os.WriteFile(filepath.Join(postDir, "a.md"), []byte(""), 0644)
		os.WriteFile(filepath.Join(postDir, "b.md"), []byte(""), 0644)

		_, err := readPost(postDir, "multi-md")
		if err == nil {
			t.Error("expected error for multiple markdown files")
		}
	})

	t.Run("MultipleJSON", func(t *testing.T) {
		postDir := filepath.Join(tempDir, "multi-json")
		os.Mkdir(postDir, 0755)
		os.WriteFile(filepath.Join(postDir, "test.md"), []byte(""), 0644)
		os.WriteFile(filepath.Join(postDir, "a.json"), []byte("{}"), 0644)
		os.WriteFile(filepath.Join(postDir, "b.json"), []byte("{}"), 0644)

		_, err := readPost(postDir, "multi-json")
		if err == nil {
			t.Error("expected error for multiple json files")
		}
	})
}

func TestScanDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "scantest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a few post directories
	posts := []struct {
		slug string
		date string
	}{
		{"old-post", "2023-01-01"},
		{"new-post", "2024-01-01"},
		{"mid-post", "2023-06-01"},
	}

	for _, p := range posts {
		dir := filepath.Join(tempDir, p.slug)
		os.Mkdir(dir, 0755)
		os.WriteFile(filepath.Join(dir, "content.md"), []byte("content"), 0644)
		os.WriteFile(filepath.Join(dir, "meta.json"), []byte(`{"date": "`+p.date+`"}`), 0644)
	}

	scanned, err := ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if len(scanned) != 3 {
		t.Errorf("expected 3 posts, got %d", len(scanned))
	}

	// Verify sorting (newest first)
	if scanned[0].Slug != "new-post" {
		t.Errorf("expected first post to be new-post, got %s", scanned[0].Slug)
	}
	if scanned[2].Slug != "old-post" {
		t.Errorf("expected last post to be old-post, got %s", scanned[2].Slug)
	}
}
