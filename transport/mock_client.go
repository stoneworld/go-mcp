package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"

	"go-mcp/pkg"
)

type MockClientTransport struct {
	receiver ClientReceiver

	in  *bufio.ReadWriter
	out *bufio.ReadWriter

	cancel context.CancelFunc
}

func NewMockClientTransport(in *bufio.ReadWriter, out *bufio.ReadWriter) *MockClientTransport {
	return &MockClientTransport{}
}

func (t *MockClientTransport) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go func() {
		for {
			t.receive(ctx)
		}
	}()

	return nil
}

func (t *MockClientTransport) Send(ctx context.Context, msg Message) error {
	if _, err := t.out.Write(append(msg, "\n"...)); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	if err := t.out.Flush(); err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}
	return nil
}

func (t *MockClientTransport) SetReceiver(receiver ClientReceiver) {
	t.receiver = receiver
}

func (t *MockClientTransport) Close() error {
	t.cancel()
	return nil
}

func (t *MockClientTransport) receive(ctx context.Context) {
	defer pkg.Recover()

	line, err := t.in.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
			return
		}
		return
	}
	t.receiver.Receive(ctx, line)
}
