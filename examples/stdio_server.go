package examples

import (
	"log"

	"go-mcp/server"
	"go-mcp/transport"
)

func main() {
	stdioTransport, err := transport.NewStdioServerTransport()
	if err != nil {
		log.Fatalf("NewStdioServerTransport: %v", err)
	}
	server, err := server.NewServer(stdioTransport)
	if err != nil {
		log.Fatalf("NewMCPServer: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("NewMCPServer: %v", err)
	}
}
