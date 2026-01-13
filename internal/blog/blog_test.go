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
