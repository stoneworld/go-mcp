package transport

import (
	"context"
	"testing"
	"time"
)

func TestSSE(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var (
		err    error
		svr    ServerTransport
		client ClientTransport
	)

	if svr, err = NewSSEServerTransport("0.0.0.0:8181"); err != nil {
		t.Errorf("NewSSEServerTransport failed: %v", err)
		return
	}

	time.Sleep(time.Second)

	if client, err = NewSSEClientTransport(ctx, "http://127.0.0.1:8181/sse"); err != nil {
		t.Errorf("NewSSEClientTransport failed: %v", err)
		return
	}

	testClient2Server(t, client, svr)
}
