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
	transport transport.Transport

	reqID2respChan map[int]chan *protocol.JSONRPCResponse

	notifyMethod2handler map[protocol.Method]func(notifyParam interface{})

	requestID atomic.Int64

	logger pkg.Logger
}

func NewClient(t transport.Transport, opts ...Option) (*Client, error) {
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

func WithNotifyHandler(notifyMethod2handler map[protocol.Method]func(notifyParam interface{})) Option {
	return func(s *Client) {
		s.notifyMethod2handler = notifyMethod2handler
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
