package main

import (
	"log"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/zhaojunlucky/ptolemaios/pkg/tool/macos"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Macbook Pro mcp tool 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool handler
	macos.Init(s)

	log.Println("Starting MCP server on :8080/mcp")
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithHeartbeatInterval(30*time.Second),
		server.WithStateLess(false),
	)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
