package transport

import (
	"context"
	"testing"
)

func TestSSE(t *testing.T) {
	var (
		err    error
		svr    ServerTransport
		client ClientTransport
	)

	if svr, err = NewSSEServerTransport("0.0.0.0:8181"); err != nil {
		t.Fatalf("NewSSEServerTransport failed: %v", err)
	}

	if client, err = NewSSEClientTransport(context.Background(), "http://127.0.0.1:8181/sse"); err != nil {
		t.Fatalf("NewSSEClientTransport failed: %v", err)
	}

	testTransport(t, client, svr)
}
