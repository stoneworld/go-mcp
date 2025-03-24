package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"

	"go-mcp/pkg"
)

const mockSessionID = "mock"

type MockServerTransport struct {
	receiver ServerReceiver

	in  *bufio.ReadWriter
	out *bufio.ReadWriter

	cancel context.CancelFunc
}

func NewMockServerTransport(in *bufio.ReadWriter, out *bufio.ReadWriter) *MockServerTransport {
	return &MockServerTransport{}
}

func (t *MockServerTransport) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go func() {
		for {
			t.receive(ctx)
		}
	}()

	return nil
}

func (t *MockServerTransport) Send(ctx context.Context, sessionID string, msg Message) error {
	if _, err := t.out.Write(append(msg, "\n"...)); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	if err := t.out.Flush(); err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}
	return nil
}

func (t *MockServerTransport) SetReceiver(receiver ServerReceiver) {
	t.receiver = receiver
}

func (t *MockServerTransport) Close() error {
	t.cancel()
	return nil
}

func (t *MockServerTransport) receive(ctx context.Context) {
	defer pkg.Recover()
	line, err := t.in.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
			return
		}
		return
	}
	t.receiver.Receive(ctx, mockSessionID, line)
}
