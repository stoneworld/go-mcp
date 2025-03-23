package transport

import (
	"context"
)

type sseClientTransport struct {
	receiver ClientReceiver
}

func NewSSEClientTransport() (ClientTransport, error) {
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

func (c *sseClientTransport) SetReceiver(receiver ClientReceiver) {
	c.receiver = receiver
}

func (c *sseClientTransport) Close() error {
	// TODO implement me
	panic("implement me")
}
