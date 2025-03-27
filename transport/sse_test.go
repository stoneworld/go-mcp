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
		t.Errorf("NewSSEServerTransport failed: %v", err)
		return
	}

	if client, err = NewSSEClientTransport(context.Background(), "http://127.0.0.1:8181/sse"); err != nil {
		t.Errorf("NewSSEClientTransport failed: %v", err)
		return
	}

	testTransport(t, client, svr)
}
