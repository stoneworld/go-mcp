package examples

import (
	"context"
	"log"

	"go-mcp/client"
	"go-mcp/transport"
)

func main() {
	stdioTransport, err := transport.NewStdioClientTransport("")
	if err != nil {
		log.Fatalf("NewStdioClientTransport: %v", err)
	}
	client, err := client.NewClient(stdioTransport)
	if err != nil {
		log.Fatalf("NewMCPClient: %v", err)
	}

	client.CallTool(context.Background())
}
