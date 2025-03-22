package examples

import (
	"log"

	"go-mcp/server"
	"go-mcp/transport"
)

func main() {
	sseTransport, err := transport.NewSSEServerTransport()
	if err != nil {
		log.Fatalf("NewSSEServerTransport: %v", err)
	}
	server, err := server.NewServer(sseTransport)
	if err != nil {
		log.Fatalf("NewMCPServer: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("NewMCPServer: %v", err)
	}
}
