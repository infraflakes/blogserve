package blog

import (
	"encoding/json"
	"fmt"
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
		if err != nil {
			// Skip invalid posts but log the error or handle it as needed
			continue
		}
		posts = append(posts, post)
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

	entries, err := os.ReadDir(dir)
	if err != nil {
		return post, err
	}

	var mdFiles []string
	var jsonFiles []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext == ".md" {
			mdFiles = append(mdFiles, entry.Name())
		} else if ext == ".json" {
			jsonFiles = append(jsonFiles, entry.Name())
		}
	}

	// Enforce exactly one markdown file
	if len(mdFiles) == 0 {
		return post, fmt.Errorf("no markdown file found in %s", dir)
	}
	if len(mdFiles) > 1 {
		return post, fmt.Errorf("multiple markdown files found in %s: %v", dir, mdFiles)
	}

	// Enforce at most one metadata file
	if len(jsonFiles) > 1 {
		return post, fmt.Errorf("multiple json files found in %s: %v", dir, jsonFiles)
	}

	// Read markdown content
	content, err := os.ReadFile(filepath.Join(dir, mdFiles[0]))
	if err != nil {
		return post, err
	}
	post.Content = string(content)

	// Read metadata if exists
	if len(jsonFiles) == 1 {
		metaFile, err := os.ReadFile(filepath.Join(dir, jsonFiles[0]))
		if err == nil {
			var meta PostMetadata
			if err := json.Unmarshal(metaFile, &meta); err == nil {
				post.Metadata = meta
			}
		}
	}

	return post, nil
}
