package main

import (
	"blogserve/internal/cmd"
	"blogserve/internal/server"
	"embed"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	server.FrontendFS = frontendFS
	cmd.Execute()
}
