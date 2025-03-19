package transport

import (
	"context"
)

type sseClientTransport struct {
	receiver receiver
}

func NewSSEClientTransport() (Transport, error) {
	return &sseClientTransport{}, nil
}

func (c *sseClientTransport) Start() error {
	// TODO implement me
	panic("implement me")
}

func (c *sseClientTransport) Send(ctx context.Context, msg Message) error {
	// TODO implement me
	panic("implement me")
}

func (c *sseClientTransport) SetReceiver(receiver receiver) {
	c.receiver = receiver
}

func (c *sseClientTransport) Close() error {
	// TODO implement me
	panic("implement me")
}
