package examples

import (
	"context"
	"log"

	"go-mcp/client"
	"go-mcp/transport"
)

func main() {
	sseTransport, err := transport.NewSSEClientTransport()
	if err != nil {
		log.Fatalf("NewSSEClientTransport: %v", err)
	}
	client, err := client.NewClient(sseTransport)
	if err != nil {
		log.Fatalf("NewMCPClient: %v", err)
	}

	client.CallTool(context.Background())
}
