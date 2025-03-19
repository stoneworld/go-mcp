package transport

import "context"

type sseServerTransport struct {
	receiver receiver
}

func NewSSEServerTransport() (Transport, error) {
	return &sseServerTransport{}, nil
}

func (c *sseServerTransport) Start() error {
	// TODO implement me
	panic("implement me")
}

func (c *sseServerTransport) Send(ctx context.Context, msg Message) error {
	// TODO implement me
	panic("implement me")
}

func (c *sseServerTransport) SetReceiver(receiver receiver) {
	c.receiver = receiver
}

func (c *sseServerTransport) Close() error {
	// TODO implement me
	panic("implement me")
}
