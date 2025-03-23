package client

import (
	"context"
	"fmt"
	"sync/atomic"

	"go-mcp/pkg"
	"go-mcp/protocol"
	"go-mcp/transport"
)

type Client struct {
	transport transport.ClientTransport

	reqID2respChan map[int]chan *protocol.JSONRPCResponse

	listRootsHandler            func(ctx context.Context, request *protocol.ListRootsRequest) (*protocol.ListRootsResult, error)
	createMessagesSampleHandler func(ctx context.Context, request *protocol.CreateMessageRequest) (*protocol.CreateMessageResult, error)

	cancelledNotifyHandler func(ctx context.Context, notifyParam *protocol.CancelledNotification) error

	requestID atomic.Int64

	logger pkg.Logger
}

func NewClient(t transport.ClientTransport, opts ...Option) (*Client, error) {
	client := &Client{
		transport: t,
		logger:    &pkg.Log{},
	}
	t.SetReceiver(client)

	for _, opt := range opts {
		opt(client)
	}

	if err := client.transport.Start(); err != nil {
		return nil, fmt.Errorf("init mcp client transpor start fail: %w", err)
	}

	if err := client.initialization(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

type Option func(*Client)

func WithListRootsHandlerHandler(handler func(ctx context.Context, request *protocol.ListRootsRequest) (*protocol.ListRootsResult, error)) Option {
	return func(s *Client) {
		s.listRootsHandler = handler
	}
}

func WithCreateMessagesSampleHandler(handler func(ctx context.Context, request *protocol.CreateMessageRequest) (*protocol.CreateMessageResult, error)) Option {
	return func(s *Client) {
		s.createMessagesSampleHandler = handler
	}
}

func WithCancelNotifyHandler(handler func(ctx context.Context, notifyParam *protocol.CancelledNotification) error) Option {
	return func(s *Client) {
		s.cancelledNotifyHandler = handler
	}
}

func WithLogger(logger pkg.Logger) Option {
	return func(s *Client) {
		s.logger = logger
	}
}

func (client *Client) Close() error {
	// TODO 还有一些其他处理操作也可以放在这里
	// client.transport.Close()
	return nil
}
