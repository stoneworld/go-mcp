package transport

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type serverReceive func(ctx context.Context, sessionID string, msg []byte) error

func (r serverReceive) Receive(ctx context.Context, sessionID string, msg []byte) error {
	return r(ctx, sessionID, msg)
}

type clientReceive func(ctx context.Context, msg []byte) error

func (r clientReceive) Receive(ctx context.Context, msg []byte) error {
	return r(ctx, msg)
}

func testClient2Server(t *testing.T, client ClientTransport, server ServerTransport) {
	msg := "hello"
	expectedMsg := ""

	testSendReceive(t, client, server, func() {
		server.SetReceiver(serverReceive(func(ctx context.Context, sessionID string, msg []byte) error {
			expectedMsg = string(msg)
			return nil
		}))
	}, func() {
		if err := client.Send(context.Background(), Message(msg)); err != nil {
			t.Fatalf("client.Send() failed: %v", err)
		}
	})

	assert.Equal(t, expectedMsg, msg)
}

func testServer2Client(t *testing.T, client ClientTransport, server ServerTransport) {
	var (
		msg         = "hello"
		expectedMsg = ""
	)

	testSendReceive(t, client, server, func() {
		client.SetReceiver(clientReceive(func(ctx context.Context, msg []byte) error {
			expectedMsg = string(msg)
			return nil
		}))
	}, func() {
		// TODO： 这里需要解决获取不到sessionID的问题
		if err := server.Send(context.Background(), "$test$", Message(msg)); err != nil {
			t.Fatalf("server.Send() failed: %v", err)
		}
	})

	assert.Equal(t, expectedMsg, msg)
}

func testSendReceive(t *testing.T, client ClientTransport, server ServerTransport, setReceive, send func()) {
	setReceive()

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

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := server.Shutdown(ctx, ctx); err != nil {
			t.Errorf("server.Shutdown() failed: %v", err)
		}
	}()

	if err := client.Start(); err != nil {
		t.Fatalf("client.Start() failed: %v", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("client.Close() failed: %v", err)
		}
	}()

	send()

	time.Sleep(time.Second * 1)
}
