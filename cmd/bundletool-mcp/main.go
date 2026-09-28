package main

import (
	"context"
	"log"

	"bundletool-mcp/internal/bundletool"
	"bundletool-mcp/internal/executor"
	"bundletool-mcp/internal/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Version is injected at build time using -ldflags
var Version = "dev"

func main() {
	// 1. Dependency Injection Setup
	exec := executor.NewDefaultExecutor()

	// Ensure bundletool is downloaded and get its absolute path
	jarPath, err := bundletool.EnsureBundleTool()
	if err != nil {
		log.Fatalf("Failed to setup bundletool: %v", err)
	}

	// Initialize our core business logic client
	btClient := bundletool.NewClient(exec, jarPath)

	// 2. Initialize the MCP Server Layer
	server := mcp.NewServer(&mcp.Implementation{Name: "bundletool-mcp", Version: Version}, nil)

	// Register the tool handlers
	handlers := tools.NewHandlers(btClient)
	handlers.Register(server)

	// 4. Start Transport (Stdio)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
