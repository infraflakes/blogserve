// Package cmd implements the command-line interface for blogserve.
package cmd

import (
	"blogserve/internal/blog"
	"blogserve/internal/server"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	blogDir string
	port    int
)

var rootCmd = &cobra.Command{
	Use:   "blogserve",
	Short: "blogserve is a simple markdown-based blog engine",
	Long:  `A fast and simple blog engine that serves markdown files with a Svelte frontend.`,
	Run: func(cmd *cobra.Command, args []string) {
		cleanDir := filepath.Clean(blogDir)
		s := server.NewServer(cleanDir, port)

		// Start watcher in background
		go blog.WatchDirectory(cleanDir, func() {
			s.TriggerReload()
		})

		if err := s.Start(); err != nil {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Cobra execution error", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&blogDir, "directory", "d", ".", "directory containing blog posts")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port to serve the blog on")
}
