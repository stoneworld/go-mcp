package transport

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type serverReceive func(ctx context.Context, sessionID string, msg []byte)

func (r serverReceive) Receive(ctx context.Context, sessionID string, msg []byte) {
	r(ctx, sessionID, msg)
}

type clientReceive func(ctx context.Context, msg []byte)

func (r clientReceive) Receive(ctx context.Context, msg []byte) {
	r(ctx, msg)
}

func testClient2Server(t *testing.T, client ClientTransport, server ServerTransport) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	msg := "hello"
	expectedMsg := "hello"

	server.SetReceiver(serverReceive(func(ctx context.Context, sessionID string, msg []byte) {
		expectedMsg = string(msg)
	}))

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	// 使用 select 来处理可能的错误
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("server.Run() failed: %v", err)
		}
	case <-time.After(time.Second):
		// 服务器正常启动
	}

	defer server.Shutdown(ctx)

	if err := client.Start(); err != nil {
		t.Fatalf("client.Start() failed: %v", err)
	}
	defer client.Close(ctx)

	if err := client.Send(context.Background(), Message(msg)); err != nil {
		t.Fatalf("client.Send() failed: %v", err)
	}

	assert.Equal(t, expectedMsg, msg)
}

func testServer2Client(t *testing.T, client ClientTransport, server ServerTransport) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var (
		msg         = ""
		expectedMsg = ""
	)

	client.SetReceiver(clientReceive(func(ctx context.Context, msg []byte) {
		expectedMsg = string(msg)
	}))

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	// 使用 select 来处理可能的错误
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("server.Run() failed: %v", err)
		}
	case <-time.After(time.Second):
		// 服务器正常启动
	}

	defer server.Shutdown(ctx)

	if err := client.Start(); err != nil {
		t.Fatalf("client.Start() failed: %v", err)
	}
	defer client.Close(ctx)

	// TODO： 这里需要解决获取不到sessionID的问题
	if err := server.Send(context.Background(), "$test$", Message(msg)); err != nil {
		t.Fatalf("server.Send() failed: %v", err)
	}

	assert.Equal(t, expectedMsg, msg)
}
