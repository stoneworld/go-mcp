package client

import (
	"context"
	"fmt"
	"sync/atomic"

	"go-mcp/pkg"
	"go-mcp/protocol"
	"go-mcp/transport"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type Client struct {
	transport transport.ClientTransport

	reqID2respChan cmap.ConcurrentMap[string, chan *protocol.JSONRPCResponse]

	roots []protocol.Root

	createMessagesSampleHandler func(ctx context.Context, request *protocol.CreateMessageRequest) (*protocol.CreateMessageResult, error)

	cancelledNotifyHandler func(ctx context.Context, notifyParam *protocol.CancelledNotification) error

	requestID atomic.Int64

	initialized atomic.Bool

	ClientInfo         *protocol.Implementation
	ClientCapabilities *protocol.ClientCapabilities

	ServerCapabilities *protocol.ServerCapabilities
	ServerInfo         *protocol.Implementation
	ServerInstructions string

	logger pkg.Logger
}

func NewClient(t transport.ClientTransport, request *protocol.InitializeRequest, opts ...Option) (*Client, error) {
	client := &Client{
		transport:      t,
		logger:         pkg.DefaultLogger,
		reqID2respChan: cmap.New[chan *protocol.JSONRPCResponse](),
	}
	t.SetReceiver(client)

	for _, opt := range opts {
		opt(client)
	}

	if err := client.transport.Start(); err != nil {
		return nil, fmt.Errorf("init mcp client transpor start fail: %w", err)
	}

	if _, err := client.initialization(context.Background(), request); err != nil {
		return nil, err
	}

	return client, nil
}

type Option func(*Client)

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
	if err := client.transport.Close(); err != nil {
		return err
	}
	return nil
}
