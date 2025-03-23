package transport

import (
	"context"
	"testing"
	"time"
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

	go server.Run()

	defer server.Close(ctx)

	client.Start()
	defer client.Close()

	client.Send(context.Background(), Message(msg))

	if msg != expectedMsg {
		t.Errorf("testClient2Server msg not as expected.\ngot  = %v\nwant = %v", expectedMsg, msg)
	}
}

func testServer2Client(t *testing.T, client ClientTransport, server ServerTransport) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	msg := ""
	expectedMsg := ""

	client.SetReceiver(clientReceive(func(ctx context.Context, msg []byte) {
		expectedMsg = string(msg)
	}))

	go server.Run()
	defer server.Close(ctx)

	time.Sleep(time.Second) // 这里需要等server start ready，不太优雅，后续需要优化

	client.Start()
	defer client.Close()

	server.Send(context.Background(), "$test$", Message(msg)) // TODO： 这里需要解决获取不到sessionID的问题

	if msg != expectedMsg {
		t.Errorf("testServer2Client msg not as expected.\ngot  = %v\nwant = %v", expectedMsg, msg)
	}
}
