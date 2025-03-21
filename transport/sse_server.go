package transport

import "context"

type sseServerTransport struct {
	receiver serverReceiver
}

func NewSSEServerTransport() (ServerTransport, error) {
	return &sseServerTransport{}, nil
}

func (c *sseServerTransport) Start() error {
	// TODO implement me
	panic("implement me")
}

func (c *sseServerTransport) Send(ctx context.Context, sessionID string, msg Message) error {
	// TODO implement me
	panic("implement me")
}

func (c *sseServerTransport) SetReceiver(receiver serverReceiver) {
	c.receiver = receiver
}

func (c *sseServerTransport) receive(ctx context.Context, sessionID string, msg []byte) {
	c.receiver.Receive(ctx, sessionID, msg)
}

func (c *sseServerTransport) Close() error {
	// TODO implement me
	panic("implement me")
}
