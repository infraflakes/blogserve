package blog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type PostMetadata struct {
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

type Post struct {
	Slug     string       `json:"slug"`
	Content  string       `json:"content"`
	Metadata PostMetadata `json:"metadata"`
}

func ScanDirectory(dir string) ([]Post, error) {
	var posts []Post

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		postDir := filepath.Join(dir, entry.Name())
		post, err := readPost(postDir, entry.Name())
		if err == nil {
			posts = append(posts, post)
		}
	}

	// Sort posts by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		// If dates are equal, sort by slug
		if posts[i].Metadata.Date == posts[j].Metadata.Date {
			return posts[i].Slug < posts[j].Slug
		}
		// Descending order for newest first
		return posts[i].Metadata.Date > posts[j].Metadata.Date
	})

	return posts, nil
}

func readPost(dir, slug string) (Post, error) {
	var post Post
	post.Slug = slug

	// Read markdown content
	mdPath := filepath.Join(dir, slug+".md")
	content, err := os.ReadFile(mdPath)
	if err != nil {
		// Try content.md if slug.md doesn't exist
		mdPath = filepath.Join(dir, "content.md")
		content, err = os.ReadFile(mdPath)
		if err != nil {
			return post, err
		}
	}
	post.Content = string(content)

	// Read metadata
	metaPath := filepath.Join(dir, "metadata.json")
	metaFile, err := os.ReadFile(metaPath)
	if err == nil {
		var meta PostMetadata
		if err := json.Unmarshal(metaFile, &meta); err == nil {
			post.Metadata = meta
		}
	}

	return post, nil
}
