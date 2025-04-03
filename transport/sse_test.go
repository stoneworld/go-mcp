package transport

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestSSE(t *testing.T) {
	var (
		err    error
		svr    ServerTransport
		client ClientTransport
	)

	// Get an available port
	port, err := getAvailablePort()
	if err != nil {
		t.Fatalf("Failed to get available port: %v", err)
	}

	serverAddr := fmt.Sprintf("127.0.0.1:%d", port)
	clientURL := fmt.Sprintf("http://%s/sse", serverAddr)

	if svr, err = NewSSEServerTransport(serverAddr); err != nil {
		t.Fatalf("NewSSEServerTransport failed: %v", err)
	}

	if client, err = NewSSEClientTransport(context.Background(), clientURL); err != nil {
		t.Fatalf("NewSSEClientTransport failed: %v", err)
	}

	testTransport(t, client, svr)
}

func TestSSEHandler(t *testing.T) {
	var (
		messageUrl = "/message"
		port       int

		err    error
		svr    ServerTransport
		client ClientTransport
	)

	// Get an available port
	port, err = getAvailablePort()
	if err != nil {
		t.Fatalf("Failed to get available port: %v", err)
	}

	serverAddr := fmt.Sprintf("http://127.0.0.1:%d", port)
	serverURL := fmt.Sprintf("%s/sse", serverAddr)

	svr, handler, err := NewSSEServerTransportAndHandler(fmt.Sprintf("%s%s", serverAddr, messageUrl))
	if err != nil {
		t.Fatalf("NewSSEServerTransport failed: %v", err)
	}

	// 设置 HTTP 路由
	http.Handle("/sse", handler.HandleSSE())
	http.Handle(messageUrl, handler.HandleMessage())

	errCh := make(chan error, 1)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Use select to handle potential errors
	select {
	case err = <-errCh:
		t.Fatalf("http.ListenAndServe() failed: %v", err)
	case <-time.After(time.Second):
		// Server started normally
	}

	if client, err = NewSSEClientTransport(context.Background(), serverURL); err != nil {
		t.Fatalf("NewSSEClientTransport failed: %v", err)
	}

	testTransport(t, client, svr)
}

// getAvailablePort returns a port that is available for use
func getAvailablePort() (int, error) {
	addr, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("failed to get available port: %v", err)
	}
	defer func() {
		if err = addr.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	port := addr.Addr().(*net.TCPAddr).Port
	return port, nil
}
