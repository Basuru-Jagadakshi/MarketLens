package main

import (
	"log"

	"marketlens-go-backend/internal/config"
	"marketlens-go-backend/internal/mcpserver"
	"marketlens-go-backend/internal/repositories"
)

// This is its own binary/process (separate from cmd/api), so an MCP
// startup failure no longer takes the customer-facing REST API down with
// it — it just exits this process.
func main() {
	config.ConnectDatabase()
	store := repositories.NewStore(config.DB)

	mcpServer := mcpserver.New(store)

	log.Println("MCP server listening on :9090/mcp")
	if err := mcpserver.StartHTTP(mcpServer, ":9090"); err != nil {
		log.Fatalf("MCP server failed: %v", err)
	}
}
