package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"go-mcp/pkg"
)

type MockClientTransport struct {
	receiver ClientReceiver

	in  io.Reader
	out io.Writer

	logger pkg.Logger

	cancel context.CancelFunc
}

func NewMockClientTransport(in io.Reader, out io.Writer) *MockClientTransport {
	return &MockClientTransport{
		in:     in,
		out:    out,
		logger: pkg.DefaultLogger,
	}
}

func (t *MockClientTransport) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go func() {
		defer pkg.Recover()

		t.receive(ctx)
	}()

	return nil
}

func (t *MockClientTransport) Send(ctx context.Context, msg Message) error {
	if _, err := t.out.Write(append(msg, "\n"...)); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

func (t *MockClientTransport) SetReceiver(receiver ClientReceiver) {
	t.receiver = receiver
}

func (t *MockClientTransport) Close(ctx context.Context) error {
	t.cancel()
	return nil
}

func (t *MockClientTransport) receive(ctx context.Context) {
	s := bufio.NewScanner(t.in)

	for s.Scan() {
		if err := t.receiver.Receive(ctx, s.Bytes()); err != nil {
			t.logger.Errorf("receiver failed: %v", err)
			return
		}
	}

	if err := s.Err(); err != nil {
		if err != io.EOF {
			t.logger.Errorf("unexpected error reading input: %v", err)
		}
		return
	}
}
