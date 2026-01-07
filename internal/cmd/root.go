package cmd

import (
	"blogserve/internal/blog"
	"blogserve/internal/server"
	"fmt"
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
			fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&blogDir, "directory", "d", ".", "directory containing blog posts")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port to serve the blog on")
}
